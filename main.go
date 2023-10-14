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

	if *projectName == "" {
		fmt.Println("Please provide a project name using the -name flag.")
		os.Exit(1)
	}

	// Directory structures
	createDirs(*projectName, []string{
		"frontend",
		"frontend/static",
		"frontend/static/css",
		"frontend/static/js",
	})

	// Copy and replace placeholders in templates
	copyTemplate("main.go", "main.go", *projectName)
	copyTemplate("package.json", "frontend/package.json", *projectName)
	copyTemplate("postcss.config.js", "frontend/postcss.config.js", *projectName)
	copyTemplate("main.css", "frontend/static/css/main.css", *projectName)
	copyTemplate("tailwind.config.js", "frontend/tailwind.config.js", *projectName)
	copyTemplate("dev-setup-linux.sh", "dev-setup-macos.sh", *projectName)
	copyTemplate("dev-setup-macos.sh", "dev-setup-macos.sh", *projectName)
	copyTemplate("dev-setup-windows.sh", "dev-setup-macos.sh", *projectName)
	copyTemplate("Dockerfile.prod", "Dockerfile.prod", *projectName)
	copyTemplate("gen-feature.sh", "gen-feature.sh", *projectName)
	copyTemplate(".gitignore", ".gitignore", *projectName)
	copyTemplate(".example.env", ".example.env", *projectName)
	copyTemplate(".air.toml", ".air.toml", *projectName)

	fmt.Printf("Project '%s' has been generated!\n", *projectName)
}

func createDirs(projectName string, dirs []string) {
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(projectName, dir), os.ModePerm)
	}
}

func copyTemplate(src, dest, projectName string) {
	data, _ := templates.ReadFile("templates/" + src)
	content := strings.ReplaceAll(string(data), "{{PROJECT_NAME}}", projectName)
	os.WriteFile(filepath.Join(projectName, dest), []byte(content), os.ModePerm)
}
