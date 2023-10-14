package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templates embed.FS

func main() {
	// Command-line arguments
	projectName := flag.String("name", "", "Name of the project to be generated")
	flag.Parse()

	if *projectName != "" {
		generateProject(projectName)
		return
	}

	featureName := flag.String("gen-feature", "", "Generate a new feature")
	flag.Parse()

	if *featureName != "" {
		generateFeature(featureName)
		return
	}

}

func generateProject(projectName *string) {

	// Directory structures
	createDirs(*projectName, []string{
		"frontend",
		"frontend/static",
		"frontend/static/css",
		"frontend/static/js",
	})

	// Copy and replace placeholders in templates
	copyProjTemplate("main.go", "main.go", *projectName)
	copyProjTemplate("package.json", "frontend/package.json", *projectName)
	copyProjTemplate("postcss.config.js", "frontend/postcss.config.js", *projectName)
	copyProjTemplate("main.css", "frontend/static/css/main.css", *projectName)
	copyProjTemplate("tailwind.config.js", "frontend/tailwind.config.js", *projectName)
	copyProjTemplate("dev-setup-linux.sh", "dev-setup-macos.sh", *projectName)
	copyProjTemplate("dev-setup-macos.sh", "dev-setup-macos.sh", *projectName)
	copyProjTemplate("dev-setup-windows.sh", "dev-setup-macos.sh", *projectName)
	copyProjTemplate("Dockerfile.prod", "Dockerfile.prod", *projectName)
	copyProjTemplate("gen-feature.sh", "gen-feature.sh", *projectName)
	copyProjTemplate(".gitignore", ".gitignore", *projectName)
	copyProjTemplate(".example.env", ".example.env", *projectName)
	copyProjTemplate(".air.toml", ".air.toml", *projectName)

	fmt.Printf("Project '%s' has been generated!\n", *projectName)
}

func generateFeature(featureName *string) {
	// Directory structures
	createDirs(*featureName, []string{
		"static",
		"static/js",
		"static/css",
		"templates",
	})

	// Copy and replace placeholders in templates
	copyFeatureTemplate("handler.go.template", "handler.go", *featureName)
	copyFeatureTemplate("routes.go.template", "routes.go", *featureName)
	copyFeatureTemplate("view.go.template", "view.go", *featureName)
	copyFeatureTemplate("create.html", "templates/create.html", *featureName)
	copyFeatureTemplate("read.html", "templates/read.html", *featureName)
	copyFeatureTemplate("update.html", "templates/update.html", *featureName)
	copyFeatureTemplate("delete.html", "templates/delete.html", *featureName)
	copyFeatureTemplate("list.html", "templates/list.html", *featureName)

	fmt.Printf("Feature '%s' has been generated!\n", featureName)
	fmt.Printf("Add the following to the main.go file:\n")
	fmt.Printf("import \"%s\"\n", *featureName)
	fmt.Printf("%s.SetupRoutes()\n", *featureName)

}
func createDirs(projectName string, dirs []string) {
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(projectName, dir), os.ModePerm)
	}
}

func copyProjTemplate(src, dest, projectName string) {
	data, _ := templates.ReadFile("templates/" + src)
	content := strings.ReplaceAll(string(data), "{{PROJECT_NAME}}", projectName)
	os.WriteFile(filepath.Join(projectName, dest), []byte(content), os.ModePerm)
}

func copyFeatureTemplate(src, dest, featureName string) {
	data, _ := templates.ReadFile("templates/feature_templates" + src)
	content := strings.ReplaceAll(string(data), "{{FEATURE_NAME}}", featureName)
	os.WriteFile(filepath.Join(featureName, dest), []byte(content), os.ModePerm)
}
