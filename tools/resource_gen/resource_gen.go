package main

import (
	"embed"
	"html/template"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

//go:embed resource.go.tmpl
var resourceTemplate embed.FS

type CodeGenResource struct {
	Provider string
	Resource string
	Model    string
}

func main() {
	templateContent, err := resourceTemplate.ReadFile("resource.go.tmpl")
	if err != nil {
		log.Fatalf("Failed to read resource template file: %v", err)
	}

	providerPrompt := promptui.Prompt{
		Label: "Enter the provider",
	}

	provider, err := providerPrompt.Run()
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	resourcePrompt := promptui.Prompt{
		Label: "Enter the resource",
	}

	resource, err := resourcePrompt.Run()
	if err != nil {
		log.Fatalf("Failed to get resource: %v", err)
	}

	modelPrompt := promptui.Prompt{
		Label: "Enter the model",
	}

	model, err := modelPrompt.Run()
	if err != nil {
		log.Fatalf("Failed to get model: %v", err)
	}

	codegen := CodeGenResource{
		Provider: provider,
		Resource: resource,
		Model:    model,
	}

	tmpl, err := template.New("codegen").Parse(string(templateContent))
	if err != nil {
		log.Fatalf("Failed to parse the template: %v", err)
	}

	f, err := os.Create("resource_" + codegen.Resource + ".go")
	if err != nil {
		log.Fatalf("Failed to create resource.go file")
	}
	defer f.Close()

	err = tmpl.Execute(f, codegen)
	if err != nil {
		log.Fatalf("Failed to generate code from template: %v", err)
	}
}
