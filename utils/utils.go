package utils

import (
	"embed"
	"fmt"
	"github.com/joho/godotenv"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownloadFile(url string, dest string, projectName string) error {
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
func LoadEmbeddedEnv(envFile embed.FS) (map[string]string, error) {
	// Read the embedded .env file
	env, err := envFile.ReadFile(".env") // make sure the path is correct relative to the embedding directive
	if err != nil {
		return nil, err
	}

	// Parse the environment variables from the byte content
	return godotenv.Unmarshal(string(env))
}

type FieldInfo struct {
	Name string
	Type string
	Tag  []struct{ Name, Value string }
}

func GetStructFields(filePath string, structName string) ([]FieldInfo, error) {
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
						tag := field.Tag
						if tag != nil {
							tagValue := tag.Value[1 : len(tag.Value)-1] // Remove surrounding backticks
							var tagInfo []struct{ Name, Value string }
							var re = regexp.MustCompile(`([a-zA-Z0-9]+):"([^"]+)"`)
							matches := re.FindAllStringSubmatch(tagValue, -1)
							for _, match := range matches {
								tagInfo = append(tagInfo, struct{ Name, Value string }{match[1], match[2]})
							}
							fields = append(fields, FieldInfo{
								Name: field.Names[0].Name,
								Type: fieldType,
								Tag:  tagInfo,
							})
							continue
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

func GenerateStructWithTags(name string, fields []FieldInfo, tagNames []string) string {
	var sb strings.Builder
	ns := "\n"
	ns = fmt.Sprintf("type %s struct{\n", name)
	sb.WriteString(ns)
	for _, field := range fields {
		sb.WriteString("\t" + field.Name + " " + field.Type + " `")
		for _, tag := range field.Tag {
			if tag.Name == "form" && (tag.Value == "createdat" || tag.Value == "updatedat") {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s:\"%s\" ", tag.Name, tag.Value))
		}
		for _, tagName := range tagNames {
			tagValue := strings.ToLower(field.Name)
			if tagName == "form" && (tagValue == "createdat" || tagValue == "updatedat") {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s:\"%s\" ", tagName, tagValue))
		}
		sb.WriteString("`\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
func SetUpEnv(env string) {
	ef := ".env"
	if env == "production" {
		ef = ".env.production"
	}
	err := godotenv.Load(ef)
	if err != nil {
		log.Fatal("Error loading env file")
	}
}
