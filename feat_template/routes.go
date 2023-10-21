package {{FEAURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"database/sql"
	"net/http"
)

var feature_router *infra.Router

func SetupRoutes(db *sql.DB, middleware ...infra.Middleware) {
	var shan = func(db *sql.DB) http.Handler {
		return http.StripPrefix("/static/{{FEAURE_NAME}}/", http.FileServer(http.Dir("./{{FEAURE_NAME}}/static/")))
	}
	routes := []infra.Route{
		{Handler: List, Path: "/{{FEAURE_NAME}}/"},
		{Handler: Create, Path: "/{{FEAURE_NAME}}/create"},
		{Handler: Read, Path: "/{{FEAURE_NAME}}/read"},
		{Handler: Update, Path: "/{{FEAURE_NAME}}/update"},
		{Handler: Delete, Path: "/{{FEAURE_NAME}}/delete"},
		{Handler: shan, Path: "/static/{{FEAURE_NAME}}/"},
	}
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(db, middleware...)
}
