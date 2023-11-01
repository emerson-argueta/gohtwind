package {{FEATURE_NAME}}

import (
	"database/sql"
	"net/http"
	"{{PROJECT_NAME}}/infra"
)

var feature_router *infra.Router

func SetupRoutes(dbs map[string]*sql.DB, middleware ...infra.Middleware) {
	handle := &Handle{dbs: dbs}
	shan := http.StripPrefix("/static/{{FEATURE_NAME}}/", http.FileServer(http.Dir("./{{FEATURE_NAME}}/static/")))
	routes := []infra.Route{
		{Method: "GET", Handler: handle.List, Path: "/{{FEATURE_NAME}}"},
		{Method: "POST", Handler: handle.Create, Path: "/{{FEATURE_NAME}}"},
		{Method: "GET", Handler: handle.Read, Path: "/{{FEATURE_NAME}}/{id}"},
		{Method: "PATCH", Handler: handle.Update, Path: "/{{FEATURE_NAME}}/{id}"},
		{Method: "DELETE", Handler: handle.Delete, Path: "/{{FEATURE_NAME}}/{id}"},
		{Method: "GET", Handler: shan.ServeHTTP, Path: "/static/{{FEATURE_NAME}}/"},
	}
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(middleware...)
}
