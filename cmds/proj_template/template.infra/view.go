package infra

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type View struct {
	templates *template.Template
	basePath  string
	fp        string
}

func NewView(basePath string, fp string) (*View, error) {
	allTemplatePaths, err := collectAllTemplatePaths(basePath, fp)
	if err != nil {
		return nil, err
	}
	templates := template.New("").Funcs(TemplateHelperFuncs)
	for _, path := range allTemplatePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		_, err = templates.New(filepath.ToSlash(path)).Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	return &View{
		basePath:  basePath,
		fp:        fp,
		templates: templates,
	}, nil
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
