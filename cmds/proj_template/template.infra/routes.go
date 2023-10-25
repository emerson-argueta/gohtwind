package infra

import (
	"database/sql"
	"net/http"
)

type Route struct {
	Handler func(dbs map[string]*sql.DB) http.Handler
	Path    string
}

type Router struct {
	routes []Route
}

type Middleware func(http.Handler) http.Handler

func NewRouter(routes []Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) SetupRoutes(dbs map[string]*sql.DB, middleware ...Middleware) {
	for _, r := range r.routes {
		h := r.Handler(dbs)
		for _, m := range middleware {
			h = http.HandlerFunc(m(h).ServeHTTP)
		}
		http.Handle(r.Path, h)
	}
}