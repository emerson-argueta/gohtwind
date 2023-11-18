package {{FEATURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"net/http"
	"path/filepath"
)

const basePath = "templates"

var fp = filepath.Join("{{FEATURE_NAME}}", "templates")
var featureView func() (*infra.View,error)

func init() {
	fv, err := infra.NewView(basePath, fp)
	featureView = func() (*infra.View,error) { return fv, err }
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	fv, err := featureView()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fv.RenderTemplate(w, tmpl, data)
}

func renderPartialTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	fv, err := featureView()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fv.RenderPartialTemplate(w, tmpl, data)
}
