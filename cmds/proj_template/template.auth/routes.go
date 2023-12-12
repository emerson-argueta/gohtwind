package auth

import (
	"database/sql"
	"net/http"
	"{{PROJECT_NAME}}/infra"
)

var feature_router *infra.Router

func SetupRoutes(dbs map[string]*sql.DB, middleware ...infra.Middleware) {
	handle := &Handle{dbs: dbs}
	shan := http.StripPrefix("/static/auth/", http.FileServer(http.Dir("./auth/static/")))

	routes := []infra.Route{
		{Method: "GET", Handler: handle.LoginGet, Path: "/auth/login"},
		{Method: "POST", Handler: handle.LoginPost, Path: "/auth/login"},
		{Method: "DELETE", Handler: handle.Logout, Path: "/auth/logout"},
		//{Method: "GET", Handler: handle.RegisterGet, Path: "auth/register"},
		//{Method: "POST", Handler: handle.RegisterPost, Path: "auth/register"},
		{Method: "GET", Handler: shan.ServeHTTP, Path: "/static/auth/"},
	}
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(middleware...)
}
