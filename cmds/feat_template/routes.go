package {{FEATURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"database/sql"
	"net/http"
)

var feature_router *infra.Router

func SetupRoutes(dbs map[string]*sql.DB, middleware ...infra.Middleware) {
	var shan = func(dbs map[string]*sql.DB) http.Handler {
		return http.StripPrefix("/static/{{FEATURE_NAME}}/", http.FileServer(http.Dir("./{{FEATURE_NAME}}/static/")))
	}
	routes := []infra.Route{
		{Handler: List, Path: "/{{FEATURE_NAME}}/"},
		{Handler: Create, Path: "/{{FEATURE_NAME}}/create"},
		{Handler: Read, Path: "/{{FEATURE_NAME}}/read"},
		{Handler: Update, Path: "/{{FEATURE_NAME}}/update"},
		{Handler: Delete, Path: "/{{FEATURE_NAME}}/delete"},
		{Handler: shan, Path: "/static/{{FEATURE_NAME}}/"},
	}
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(dbs, middleware...)
}
