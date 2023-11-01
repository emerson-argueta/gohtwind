package infra

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type View struct {
	templates *template.Template
	basePath  string
	fp        string
}

func iterMap(index int, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Index": index,
		"Value": value,
	}
}
func dictFunc(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
func sliceFunc(values ...interface{}) []interface{} {
	return values
}

func NewView(basePath string, fp string) *View {
	allTemplatePaths, err := collectAllTemplatePaths(basePath, fp)
	if err != nil {
		panic(err)
	}
	templates := template.New("").Funcs(template.FuncMap{"iterMap": iterMap, "dict": dictFunc, "slice": sliceFunc})
	for _, path := range allTemplatePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		_, err = templates.New(filepath.ToSlash(path)).Parse(string(content))
		if err != nil {
			panic(err)
		}
	}

	return &View{
		basePath:  basePath,
		fp:        fp,
		templates: templates,
	}
}

func collectAllTemplatePaths(paths ...string) ([]string, error) {
	var allTemplatePaths []string
	for _, path := range paths {
		templatePaths, err := collectTemplatePaths(path, ".html")
		if err != nil {
			return nil, err
		}
		allTemplatePaths = append(allTemplatePaths, templatePaths...)
	}
	return allTemplatePaths, nil
}

func collectTemplatePaths(root string, ext string) ([]string, error) {
	var templatePaths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			templatePath := filepath.ToSlash(path)
			templatePaths = append(templatePaths, templatePath)
		}
		return nil
	})
	return templatePaths, err
}

func (v *View) RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := v.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *View) RenderPartialTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Get the specific template entry (without the base layout)
	partial, err := v.templates.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	partial, err = partial.New("").ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = partial.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
