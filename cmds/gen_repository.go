package cmds

import (
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func genRepoUsageString() string {
	return `
Usage: gohtwind gen-repository [options]
    Options:
		-feature-name string
			Name of the feature the repository is for
		-model-name string
			Name of the model the repository is for
		-db-name string
			Name of the database the model is in
		-schema-name string
			Name of the schema the model is in (postgres adapter only)
		-adapter string
			Database adapter (mysql, postgres)
	Info:
		* Generates a repository file for the specified feature and model
		* The repository file contains boilerplate code for basic CRUD operations
		* The repository file is used by the feature's handler to interact with the database

`
}

//go:embed repo_template/repo_partial.go
var repoPartialFile embed.FS

type repo struct {
	flagSet     *flag.FlagSet
	featName    *string
	modelName   *string
	dbName      *string
	schema      *string
	adapter     *string
	projectPath string
}

func newRepo() *repo {
	fgs := flag.NewFlagSet("gohtwind gen-repository", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return &repo{
		flagSet:     fgs,
		featName:    fgs.String("feature-name", "", "Name of the feature the repository is for"),
		modelName:   fgs.String("model-name", "", "Name of the model the repository is for"),
		dbName:      fgs.String("db-name", "", "Name of the database the model is in"),
		schema:      fgs.String("schema-name", "", "Name of the schema the model is in (postgres adapter only)"),
		adapter:     fgs.String("adapter", "", "Database adapter (mysql, postgres)"),
		projectPath: pjp,
	}
}
func GenRepository() {
	r := newRepo()
	args := os.Args[2:]
	err := r.flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if *r.featName == "" || *r.modelName == "" || *r.dbName == "" || *r.adapter == "" {
		fmt.Println(genRepoUsageString())
		os.Exit(1)
	}
	if *r.adapter == "postgres" && *r.schema == "" {
		fmt.Println(genRepoUsageString())
		os.Exit(1)
	}
	r.genRepoFile()
	fmt.Printf("Repository has been generated for feature %s, with model: %s!\n", *r.featName, *r.modelName)
}

func (r *repo) genRepoFile() {
	featPath := filepath.Join(r.projectPath, *r.featName)
	f := fmt.Sprintf("%s_%s_%s_repo.go",
		strings.ToLower(*r.adapter),
		strings.ToLower(*r.dbName),
		strings.ToLower(*r.modelName),
	)
	repoFile, err := os.Create(filepath.Join(featPath, f))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer repoFile.Close()
	r.writeImports(repoFile)
	r.writeDto(repoFile)
	r.writePartial(repoFile)
}

func (r *repo) writeImports(repoFile *os.File) {
	ps := strings.Split(r.projectPath, "/")
	projName := ps[len(ps)-1]
	imports := fmt.Sprintf("package %s\n\n", *r.featName)
	imports = fmt.Sprintf("%simport(\n", imports)
	imports = fmt.Sprintf("%s\t\"database/sql\"\n", imports)
	imports = fmt.Sprintf("%s\t\"time\"\n", imports)
	imports = fmt.Sprintf("%s\t\"log\"\n", imports)
	imports = fmt.Sprintf("%s\t\"%s/infra\"\n", imports, projName)
	if *r.adapter == "postgres" {
		imports = fmt.Sprintf("%s\t. \"github.com/go-jet/jet/v2/postgres\"\n", imports)
		imports = fmt.Sprintf("%s\t\"%s/.gen/%s/%s/model\"\n", imports, projName, *r.dbName, *r.schema)
		imports = fmt.Sprintf("%s\t. jet \"%s/.gen/%s/%s/table\"\n", imports, projName, *r.dbName, *r.schema)
	}
	if *r.adapter == "mysql" {
		imports = fmt.Sprintf("%s\t. \"github.com/go-jet/jet/v2/mysql\"\n", imports)
		imports = fmt.Sprintf("%s\t\"%s/.gen/%s/model\"\n", imports, projName, *r.dbName)
		imports = fmt.Sprintf("%s\t. \"%s/.gen/%s/table\"\n", imports, projName, *r.dbName)
	}
	imports = fmt.Sprintf("%s)\n\n", imports)
	_, err := repoFile.WriteString(imports)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (r *repo) writeDto(repoFile *os.File) {
	fp := filepath.Join(r.projectPath, ".gen", *r.dbName, "model", fmt.Sprintf("%s.go", *r.modelName))
	fields, err := getStructFields(fp, *r.modelName)
	fmt.Println(fields)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tags := []string{"form", "json"}
	ns := fmt.Sprintf("%sDto", *r.modelName)
	newStruct := generateStructWithTags(ns, fields, tags)
	fmt.Println(newStruct)
	_, err = repoFile.WriteString(newStruct)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (r *repo) writePartial(repoFile *os.File) {
	repoPartial, err := repoPartialFile.ReadFile("repo_template/repo_partial.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	content := string(repoPartial)
	content = strings.ReplaceAll(content, "{{DB_NAME}}", *r.dbName)
	content = strings.ReplaceAll(content, "{{MODEL_NAME}}", *r.modelName)
	_, err = repoFile.WriteString(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type FieldInfo struct {
	Name string
	Type string
}

func getStructFields(filePath string, structName string) ([]FieldInfo, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var fields []FieldInfo

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if typeSpec.Name.Name == structName {
				for _, field := range structType.Fields.List {
					if len(field.Names) > 0 {
						fieldType := ""
						switch t := field.Type.(type) {
						case *ast.SelectorExpr:
							fieldType = t.X.(*ast.Ident).Name + "." + t.Sel.Name
						default:
							fieldType = fmt.Sprint(field.Type)
						}

						fields = append(fields, FieldInfo{
							Name: field.Names[0].Name,
							Type: fieldType,
						})
					}
				}
			}
		}
	}

	return fields, nil
}

func generateStructWithTags(name string, fields []FieldInfo, tagNames []string) string {
	var sb strings.Builder
	ns := fmt.Sprintf("type %s struct{\n", name)
	sb.WriteString(ns)
	for _, field := range fields {
		sb.WriteString("\t" + field.Name + " " + field.Type + " `")
		for _, tagName := range tagNames {
			sb.WriteString(fmt.Sprintf("%s:\"%s\" ", tagName, strings.ToLower(field.Name)))
		}
		sb.WriteString("`\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
