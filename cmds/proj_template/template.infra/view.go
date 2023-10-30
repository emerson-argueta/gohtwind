package infra

import (
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

func NewView(basePath string, fp string) *View {
	allTemplatePaths, err := collectAllTemplatePaths(basePath, fp)
	if err != nil {
		panic(err)
	}
	templates := template.New("").Funcs(template.FuncMap{"iterMap": iterMap})
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
	partial, err = partial.New("").ParseFiles(filepath.Join(v.fp, tmpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = partial.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
