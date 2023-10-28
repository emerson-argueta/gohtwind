package infra

import (
	"net/http"
	"path/filepath"
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
	// Parse the base layout.
	base, err := template.ParseFiles(basePath)
	if err != nil {
		panic(err)
	}
	// Parse the feature templates and associate them with the base layout.
	templates, err := base.Funcs(template.FuncMap{"iterMap": iterMap}).ParseGlob(filepath.Join(fp, "*.gohtml"))
	if err != nil {
		panic(err)
	}
	return &View{
		basePath:  basePath,
		fp:        fp,
		templates: templates,
	}
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
