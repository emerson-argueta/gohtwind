package cmds

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func genFeatUsageString() string {
	return `
Usage: gohtwind gen-feature <feature_name>
	Info:
		* This command should not be run outside of a gohtwind project directory. 
		* Creates a module with the same name of the feature in the root of your project directory. 
		* A feature is a page within the web application. 
		* Generates boilerplate code for basic CRUD operations.
		* All the CRUD operations are done within the context of a single page. 
		* Traditionally, CRUD operations are done across multiple pages.
		* Gohtwind's approach is to keep all the CRUD operations within a single page using dialog modals. 
`
}

func GenFeature() {
	if len(os.Args) < 3 {
		fmt.Println(genFeatUsageString())
		os.Exit(1)
	}
	featureName := os.Args[2]
	if featureName == "" {
		fmt.Println(genFeatUsageString())
		os.Exit(1)
	}
	copyFeatTemplate(featureName)
	fmt.Printf("Feature '%s' has been generated!\n", featureName)
	fmt.Printf("Add the following to the main.go file:\n")
	fmt.Printf("import \"%s\"\n", featureName)
	fmt.Printf("%s.SetupRoutes(dbs, infra.LoggingMiddleware)\n", featureName)

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
