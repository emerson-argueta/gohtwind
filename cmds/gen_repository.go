package cmds

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func genRepoUsageString() string {
	return `
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

func GenRepository() {
	generateRepo := flag.NewFlagSet("gen-repository", flag.ExitOnError)
	repoFeatName := generateRepo.String("feature-name", "", "Name of the feature the repository is for")
	repoModelName := generateRepo.String("model-name", "", "Name of the model the repository is for")
	repoDbName := generateRepo.String("db-name", "", "Name of the database the model is in")
	repoSchema := generateRepo.String("schema-name", "", "Name of the schema the model is in (postgres adapter only)")
	repoAdapter := generateRepo.String("adapter", "", "Database adapter (mysql, postgres)")
	args := os.Args[2:]
	generateRepo.Parse(args)
	if *repoFeatName == "" || *repoModelName == "" || *repoDbName == "" || *repoAdapter == "" {
		fmt.Println(genRepoUsageString())
		os.Exit(1)
	}
	if *repoAdapter == "postgres" && *repoSchema == "" {
		fmt.Println(genRepoUsageString())
		os.Exit(1)
	}
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
