package cmds

import (
	"embed"
	"fmt"
	"gohtwind/utils"
	"io/fs"
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

func genProjUsageString() string {
	return `
Usage: gohtwind new <project_name> 
	Info:
		* Creates a new project with the name you specify.
		* The project is created in the current directory.
`
}

func GenProject() {
	if len(os.Args) < 3 {
		fmt.Println(genProjUsageString())
		os.Exit(1)
	}
	projectName := os.Args[2]
	if projectName == "" {
		fmt.Println(genProjUsageString())
		os.Exit(1)
	}
	copyProjTemplate(projectName)
	envMap, err := utils.LoadEmbeddedEnv(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	tc := NewTailwindCompiler(envMap)
	tc.downloadCompiler(projectName, "frontend/bin/tailwindcss")
	err = utils.DownloadFile("https://unpkg.com/htmx.org/dist/htmx.min.js", "frontend/static/js/htmx.min.js", projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Project '%s' has been generated!\n", projectName)
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
