package infra

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ViewTemplate struct {
	BasePath     string
	Path         string
	PartialPaths []string
}

type View struct {
	templates    *template.Template
	viewTemplate *ViewTemplate
}

func NewView(vt *ViewTemplate) (*View, error) {
	paths := append([]string{vt.BasePath, vt.Path}, vt.PartialPaths...)
	allTemplatePaths, err := collectAllTemplatePaths(paths...)
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
		templates:    templates,
		viewTemplate: vt,
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

func (v *View) renderTemplate(w http.ResponseWriter, data interface{}) {
	err := v.templates.ExecuteTemplate(w, v.viewTemplate.Path, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderTemplate(w http.ResponseWriter, vt *ViewTemplate, data interface{}) error {
	fv, err := NewView(vt)
	if err != nil {
		return err
	}
	fv.renderTemplate(w, data)
	return nil
}
