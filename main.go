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
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed proj_template/*
var projTemplate embed.FS

//go:embed feat_template/*
var featTemplate embed.FS

//go:embed .env
var envFile embed.FS

var cmdFuncs = map[string]func(){
	"new":            generateProject,
	"gen-feature":    generateFeature,
	"gen-models":     generateModels,
	"gen-repository": generateRepository,
}

func main() {
	if len(flag.Args()) == 0 {
		fmt.Println(usageString())
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	if f, ok := cmdFuncs[cmd]; !ok {
		fmt.Println(usageString())
		os.Exit(1)
	} else {
		f()
	}
}

func usageString() string {
	return `Usage: gohtwind new [options]
			Options: 
				-name string
					Name of the project to be generated
			Usage: gohtwind gen-feature [options]
			Options:
				-name string
					Name of the feature to be generated
			Usage: gohtwind gen-models [options]
			Options:
				-adapter string
					Database adapter (mysql, postgres)
				-dsn string
					Database connection string
					postgres ex: <username>:<password>@tcp(<host>:<port>)/<dbname>
					mysql ex: <username>:<password>@tcp(<host>:<port>)/<dbname
				-schema string
					Database schema (postgres adapter only)
			Usage: gohtwind gen-repository [options]
			Options:
				-feature-name string
					Name of the feature the repository is for
				-model-name string
					Name of the model the repository is for
				-db-name-or-schema string
					Name of the database (mysql) or schema (postgres) the model is in
				-adapter string
					Database adapter (mysql, postgres)
	`
}

func generateProject() {
	projectName := flag.String("name", "", "Name of the project to be generated")
	flag.Parse()
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

func generateFeature() {
	featureName := flag.String("name", "", "name of the feature to be generated")
	flag.Parse()
	copyFeatTemplate(*featureName)
	fmt.Printf("Feature '%s' has been generated!\n", *featureName)
	fmt.Printf("Add the following to the main.go file:\n")
	fmt.Printf("import \"%s\"\n", *featureName)
	fmt.Printf("%s.SetupRoutes(dbs, infra.LoggingMiddleware)\n", *featureName)

}

func generateModels() {
	genModelsFlags := flag.NewFlagSet("models", flag.ExitOnError)
	modelsAdapter := genModelsFlags.String("adapter", "", "Database adapter (mysql, postgres)")
	u := `Database connection string
			postgres ex: <username>:<password>@tcp(<host>:<port>)/<dbname>
			mysql ex: <username>:<password>@tcp(<host>:<port>)/<dbname`
	modelsDsn := genModelsFlags.String("dsn", "", u)
	modelsSchema := genModelsFlags.String("schema", "", "Database schema (postgres adapter only)")
	flag.Parse()
	dsnArg := fmt.Sprintf("-dsn=%s", *modelsDsn)
	adapterArg := fmt.Sprintf("-adapter=%s", *modelsAdapter)
	var schemaArg string
	if *modelsSchema == "" {
		schemaArg = fmt.Sprintf("-schema=%s", *modelsSchema)
		exec.Command("./bin/jet", "-path=./.gen", dsnArg, schemaArg, adapterArg).Run()
	} else {
		exec.Command("./bin/jet", "-path=./.gen", dsnArg, adapterArg).Run()
	}
}

func generateRepository() {
	generateRepo := flag.NewFlagSet("gen-repository", flag.ExitOnError)
	repoFeatName := generateRepo.String("feature-name", "", "Name of the feature the repository is for")
	repoModelName := generateRepo.String("model-name", "", "Name of the model the repository is for")
	repoDbName := generateRepo.String("db-name", "", "Name of the database the model is in")
	repoSchema := generateRepo.String("schema-name", "", "Name of the schema the model is in (postgres adapter only)")
	repoAdapter := generateRepo.String("adapter", "", "Database adapter (mysql, postgres)")
	flag.Parse()
	// find feature directory
	projPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ps := strings.Split(projPath, "/")
	projName := ps[len(ps)-1]
	featPath := filepath.Join(projPath, *repoFeatName)
	// create repository.go file in feature directory
	repoFile, err := os.Create(filepath.Join(featPath, "repository.go"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer repoFile.Close()
	// write to repository.go file
	repoFile.WriteString(fmt.Sprintf("package %s\n\n", *repoFeatName))
	// write jet sql driver import
	if *repoAdapter == "postgres" {
		repoFile.WriteString(fmt.Sprintf("import _ \"github.com/go-jet/jet/v2/postgres\"\n"))
		repoFile.WriteString(fmt.Sprintf("import \"%s/.gen/%s/%s/model\"\n\n", projName, *repoDbName, *repoSchema))
		repoFile.WriteString(fmt.Sprintf("import \"%s/.gen/%s/%s/table\"\n\n", projName, *repoDbName, *repoSchema))
	}
	if *repoAdapter == "mysql" {
		repoFile.WriteString(fmt.Sprintf("import _ \"github.com/go-jet/jet/v2/mysql\"\n\n"))
		repoFile.WriteString(fmt.Sprintf("import \"%s/.gen/%s/model\"\n\n", projName, *repoDbName))
		repoFile.WriteString(fmt.Sprintf("import \"%s/.gen/%s/table\"\n\n", projName, *repoDbName))
	}
	// TODO: write basic CRUD functions

	fmt.Printf("Repository has been generated for feature %s, with model: %s!\n", *repoFeatName, *repoModelName)
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
