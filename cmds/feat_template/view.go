package {{FEATURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, vt *infra.ViewTemplate, data interface{}) {
	fv, err := infra.NewView(vt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fv.RenderTemplate(w, data)
}