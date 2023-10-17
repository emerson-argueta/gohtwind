package {{FEATURE_NAME}}

import (
	"{{PROJECT_NAME}}/infra"
	"net/http"
)

var feature_router *infra.Router

func SetupRoutes(middleware ...infra.Middleware) {
	var shan = http.StripPrefix("/static/{{FEATURE_NAME}}/", http.FileServer(http.Dir("./{{FEATURE_NAME}}/static/")))
	routes := []infra.Route{
		{Handler: http.HandlerFunc(List), Path: "/{{FEATURE_NAME}}/"},
		{Handler: http.HandlerFunc(Create), Path: "/{{FEATURE_NAME}}/create"},
		{Handler: http.HandlerFunc(Read), Path: "/{{FEATURE_NAME}}/read"},
		{Handler: http.HandlerFunc(Update), Path: "/{{FEATURE_NAME}}/update"},
		{Handler: http.HandlerFunc(Delete), Path: "/{{FEATURE_NAME}}/delete"},
		{Handler: shan, Path: "/static/{{FEATURE_NAME}}/"},
	}
	feature_router = infra.NewRouter(routes)
	feature_router.SetupRoutes(middleware...)
}
