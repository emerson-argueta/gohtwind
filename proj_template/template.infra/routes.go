package infra

import (
	"net/http"
)

type Route struct {
	Handler http.Handler
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

func (r *Router) SetupRoutes(middleware ...Middleware) {
	for _, r := range r.routes {
		h := r.Handler
		for _, m := range middleware {
			h = http.HandlerFunc(m(h).ServeHTTP)
		}
		http.Handle(r.Path, h)
	}
}
