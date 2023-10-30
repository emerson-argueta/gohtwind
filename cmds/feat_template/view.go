package {{FEATURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"net/http"
	"path/filepath"
)

const basePath = "templates"

var fp = filepath.Join("{{FEATURE_NAME}}", "templates")
var feature_view *infra.View

func init() {
	feature_view = infra.NewView(basePath, fp)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	feature_view.RenderTemplate(w, tmpl, data)
}

func renderPartialTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	feature_view.RenderPartialTemplate(w, tmpl, data)
}
