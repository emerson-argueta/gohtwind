package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed proj_template/*
var projTemplate embed.FS

//go:embed feat_template/*
var featTemplate embed.FS

//go:embed .env
var envFile embed.FS

func main() {
	// Command-line arguments
	projectName := flag.String("name", "", "Name of the project to be generated")
	featureName := flag.String("gen-feature", "", "Generate a new feature")
	flag.Parse()

	if *projectName != "" {
		generateProject(projectName)
	} else if *featureName != "" {
		generateFeature(featureName)
	} else {
		fmt.Println("Usage: gohtwind [-name project-name] [-gen-feature feature-name]")
		os.Exit(1)
	}
}

// loadEmbeddedEnv reads the .env file from the embedded file system.
func loadEmbeddedEnv() (map[string]string, error) {
	// Read the embedded .env file
	env, err := envFile.ReadFile(".env") // make sure the path is correct relative to the embedding directive
	if err != nil {
		return nil, err
	}

	// Parse the environment variables from the byte content
	return godotenv.Unmarshal(string(env))
}

func generateProject(projectName *string) {
	copyProjTemplate(*projectName)
	envMap, err := loadEmbeddedEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	tc := NewTailwindCompiler(envMap)
	tc.downloadCompiler(*projectName, "frontend/bin/tailwindcss")
	err = downloadFile("https://unpkg.com/htmx.org/dist/htmx.min.js", "frontend/static/js/htmx.min.js", *projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Project '%s' has been generated!\n", *projectName)
}

func generateFeature(featureName *string) {
	copyFeatTemplate(*featureName)
	fmt.Printf("Feature '%s' has been generated!\n", *featureName)
	fmt.Printf("Add the following to the main.go file:\n")
	fmt.Printf("import \"%s\"\n", *featureName)
	fmt.Printf("%s.SetupRoutes(infra.LoggingMiddleware)\n", *featureName)

}

func copyProjTemplate(projectName string) {
	// Create the project directory
	projectPath := filepath.Join(".", projectName)
	os.MkdirAll(projectPath, os.ModePerm)

	// Walk through the embedded file system
	err := fs.WalkDir(projTemplate, "proj_template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Construct the target path
		targetPath := filepath.Join(projectPath, strings.TrimPrefix(path, "proj_template"))
		targetPath = strings.ReplaceAll(targetPath, "template.", "")
		if d.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		if strings.Contains(targetPath, "blank.txt") {
			return nil
		}
		// If it's a file, read and copy it
		data, err := projTemplate.ReadFile(path)
		if err != nil {
			return err
		}
		// If you need to replace placeholders within the files, do it here.
		content := strings.ReplaceAll(string(data), "{{PROJECT_NAME}}", projectName)
		// Write the file to the target directory
		return os.WriteFile(targetPath, []byte(content), os.ModePerm)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func copyFeatTemplate(featureName string) {
	// get the project name from the current directory
	projPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ps := strings.Split(projPath, "/")
	projName := ps[len(ps)-1]

	featurePath := filepath.Join(".", featureName)
	os.MkdirAll(featurePath, os.ModePerm)
	err = fs.WalkDir(featTemplate, "feat_template", func(path string, d fs.DirEntry, err error) error {
		targetPath := filepath.Join(featurePath, strings.TrimPrefix(path, "feat_template"))
		targetPath = strings.ReplaceAll(targetPath, "template.", "")
		if d.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		if strings.Contains(targetPath, "blank.txt") {
			return nil
		}
		data, err := featTemplate.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		content = strings.ReplaceAll(content, "{{FEATURE_NAME}}", featureName)
		content = strings.ReplaceAll(content, "{{PROJECT_NAME}}", projName)
		return os.WriteFile(targetPath, []byte(content), os.ModePerm)

	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func downloadFile(url string, dest string, projectName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath.Join(projectName, dest))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
