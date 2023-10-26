package {{FEATURE_NAME}}

import (
	"database/sql"
	"net/http"
	"{{PROJECT_NAME}}/infra"
)

var feature_router *infra.Router

func SetupRoutes(dbs map[string]*sql.DB, middleware ...infra.Middleware) {
	var shan = func(dbs map[string]*sql.DB) http.Handler {
		return http.StripPrefix("/static/{{FEATURE_NAME}}/", http.FileServer(http.Dir("./{{FEATURE_NAME}}/static/")))
	}
	routes := []infra.Route{
		{Method: "GET", Handler: List, Path: "/{{FEATURE_NAME}}/"},
		{Method: "POST", Handler: Create, Path: "/{{FEATURE_NAME}}/create"},
		{Method: "GET", Handler: Read, Path: "/{{FEATURE_NAME}}/read"},
		{Method: "PATCH", Handler: Update, Path: "/{{FEATURE_NAME}}/update"},
		{Method: "DELETE", Handler: Delete, Path: "/{{FEATURE_NAME}}/delete"},
		{Method: "GET", Handler: shan, Path: "/static/{{FEATURE_NAME}}/"},
	}
	feature_router = infra.NewRouter(dbs, routes)
	feature_router.SetupRoutes(middleware...)
}
