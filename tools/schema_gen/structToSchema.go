package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

//go:embed schema.tmpl
var schemaTemplate embed.FS

// //go:embed attr.tmpl
// var attrTemplate embed.FS

type StructInfo struct {
	StructName string
	Variables  []VariableInfo
}

type VariableInfo struct {
	Name string
	Type string
	Tag  string
}

type SchemaFunc struct {
	FuncName string
	SchemaList []SchemaVars
}

type SchemaVars struct {
	Name string 
	Attr string 
	Desc string 
	IsSpecial bool 
	AttrSpecial string 
}

func main() {
	var input string

	// store json key-value into dictionary map.
	dictionary, shouldReturn := getDictionary("types_conversion.json")
	if shouldReturn {
		return
	}

	// read the file to generate the schema from tfsdk tagged struct
	input, shouldReturn1 := getStructFile(input)
	if shouldReturn1 {
		return
	}

	structs := getStructs(input)
	schemaFuncs := getSchemaFuncs(structs, dictionary)
	f, err := os.Create("schema.gen")
	if err != nil {
		log.Fatalf("Failed to create schema.gen file")
	}
	defer f.Close()
	
	for _, sF := range schemaFuncs {
		// if !contains(schemaSkipList,strings.Trim(sF.FuncName,"Ome")){
			templateContent, err := schemaTemplate.ReadFile("schema.tmpl")
			if err != nil {
				log.Fatalf("Failed to read resource template file: %v", err)
			}
			tmpl, err := template.New("schemaFuncs").Parse(string(templateContent))
			if err != nil {
				log.Fatalf("Failed to parse the template: %v", err)
			}
			err = tmpl.Execute(f, sF)
			if err != nil {
				log.Fatalf("Failed to generate code from template: %v", err)
			}
			f.WriteString("\n\n")
	}
}

func getSchemaFuncs(structs []StructInfo, dictionary map[string]string) ([]SchemaFunc){
	schemaFuncs := make([]SchemaFunc,0)
	for _, structin := range structs {
		schemaList := make([]SchemaVars, 0)
		funcName := structin.StructName
		for _, vars := range structin.Variables {
			regName := regexp.MustCompile(`tfsdk:"(.*?)"`)
			matchName := regName.FindStringSubmatch(vars.Tag)
			reDis := regexp.MustCompile(`([a-z])([A-Z])`)
			description := reDis.ReplaceAllString(vars.Name, "${1} ${2}")
			attr, found := dictionary[vars.Type]
			isSpecial := false
			attrSpecial := ""
			if !found && attr == "" {
				isSpecial = true
				listReg := `^(\[\]).*$`
				re := regexp.MustCompile(listReg)
				if strings.Contains(vars.Type, "types") {
					attr = "schema.ListAttribute"
					attrSpecial = "ElementType: " + strings.Trim(vars.Type,"[]") + "Type,"
				} else if ok := re.MatchString(vars.Type); ok {
					attr = re.ReplaceAllString(vars.Type, "schema.SetNestedAttribute")
					attrSpecial = "NestedObject: schema.NestedAttributeObject{Attributes: Ome" + strings.Trim(vars.Type,"[]") + "Schema(),},"
				} else {
					attr = "schema.SingleNestedAttribute"
					attrSpecial = "Attributes: Ome" + vars.Type + "Schema(),"
				}
			}
			schemaVar := SchemaVars{
				Name: matchName[1],
				Attr: attr,
				Desc: description,
				IsSpecial: isSpecial,
				AttrSpecial: attrSpecial,
			}
			schemaList = append(schemaList, schemaVar)
		}
		schemaFunc := SchemaFunc{
			FuncName:   funcName,
			SchemaList: schemaList,
		}
		schemaFuncs = append(schemaFuncs, schemaFunc)
	}
	return schemaFuncs
}

func getStructs(input string) []StructInfo {
	structs := make([]StructInfo, 0)
	matches, ok := isStruct(input)
	if ok {
		for _, match := range matches {
			structName := match[1]
			extractedText := match[2]
			variables := make([]VariableInfo, 0)
			lines := strings.Split(extractedText, "\n")
			variables = getVariables(lines, variables)
			structInfo := StructInfo{
				StructName: structName,
				Variables:  variables,
			}
			structs = append(structs, structInfo)
		}
	}
	return structs
}

func getVariables(lines []string, variables []VariableInfo) []VariableInfo {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				variable := VariableInfo{
					Name: parts[0],
					Type: parts[1],
					Tag:  parts[2],
				}
				variables = append(variables, variable)
			}
		}
	}
	return variables
}

func isStruct(input string) ([][]string, bool) {
	regexPattern := `type\s+(\w+)\s+struct\s*{([^}]*)}`
	re := regexp.MustCompile(regexPattern)
	matches := re.FindAllStringSubmatch(input, -1)
	if len(matches) > 0 {
		return matches, true
	}
	return nil, false
}

func getStructFile(input string) (string, bool) {
	file, err := os.Open("code.tfsdk.plc")
	if err != nil {
		fmt.Println(err)
		return "", true
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return "", true
	}
	return input, false
}

func getDictionary(filepath string) (map[string]string, bool) {
	dictionary := make(map[string]string,0)
	conversion, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, true
	}
	defer conversion.Close()
	decoder := json.NewDecoder(conversion)
	err = decoder.Decode(&dictionary)
	if err != nil {
		fmt.Println(err)
		return nil, true
	}
	return dictionary, false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}