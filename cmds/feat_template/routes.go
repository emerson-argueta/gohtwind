package {{FEATURE_NAME}}

import (
	"database/sql"
	"net/http"
	"{{PROJECT_NAME}}/infra"
)

var feature_router *infra.Router

func SetupRoutes(dbs map[string]*sql.DB, middleware ...infra.Middleware) {
	handle := &Handle{dbs: dbs}
	routes := []infra.Route{
		{Method: "GET", HandleFunc: handle.List, Path: "/{{FEATURE_NAME}}"},
		{Method: "POST", HandleFunc: handle.Create, Path: "/{{FEATURE_NAME}}"},
		{Method: "GET", HandleFunc: handle.Read, Path: "/{{FEATURE_NAME}}/{id}"},
		{Method: "PATCH", HandleFunc: handle.Update, Path: "/{{FEATURE_NAME}}/{id}"},
		{Method: "DELETE", HandleFunc: handle.Delete, Path: "/{{FEATURE_NAME}}/{id}"},
	}
	http.Handle("/{{FEATURE_NAME}}/static/",http.StripPrefix("/{{FEATURE_NAME}}/static/", http.FileServer(http.Dir("./{{FEATURE_NAME}}/static/"))))
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(middleware...)
}
