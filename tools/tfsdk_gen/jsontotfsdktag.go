package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func camelCaseToUnderscore(input string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snakeCase := re.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snakeCase)
}

func main() {
	file, err := os.Open("code.struct.plc")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	conversion, err := os.Open("conversion.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conversion.Close()
	var dictionary map[string]string

	decoder := json.NewDecoder(conversion)
	err = decoder.Decode(&dictionary)
	if err != nil {
		fmt.Println(err)
		return
	}

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		reStruct := regexp.MustCompile(`type (.*?) struct`)
		matchStruct := reStruct.FindStringSubmatch(line)
		if len(matchStruct) > 1 {
			structName := matchStruct[1]
			line = strings.ReplaceAll(line, structName, "Ome"+structName)
		}
		reCamel := regexp.MustCompile(`json:"(.*?)"`)
		matchCamel := reCamel.FindStringSubmatch(line)
		if len(matchCamel) > 1 {
			camelCase := matchCamel[1]
			underscore := camelCaseToUnderscore(camelCase)
			line = strings.ReplaceAll(line, camelCase, underscore)
		}
		for key, value := range dictionary {
			line = strings.ReplaceAll(line, key, value)
		}
		lines = append(lines, line)
	}

	output := strings.Join(lines, "\n")
	outputFile, err := os.Create("code.tfsdk.gen")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	fmt.Fprintln(writer, output)
}
