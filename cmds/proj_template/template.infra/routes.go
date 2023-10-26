package infra

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	Method  string
	Handler func(dbs map[string]*sql.DB) http.Handler
	Path    string
}
type internalRoute struct {
	method  string
	handler http.Handler
	pattern *regexp.Regexp
}

type Router struct {
	routes map[string]internalRoute
}

type Middleware func(http.Handler) http.Handler

func NewRouter(dbs map[string]*sql.DB, routes []Route) *Router {
	r := &Router{
		routes: make(map[string]internalRoute),
	}

	for _, route := range routes {
		pattern := pathToRegex(route.Path)
		key := route.Method + " " + pattern.String()
		r.routes[key] = internalRoute{
			handler: route.Handler(dbs),
			pattern: pattern,
			method:  route.Method,
		}
	}

	return r
}

func pathToRegex(path string) *regexp.Regexp {
	// Directly replace the "{id}" placeholder with a regex pattern.
	pattern := strings.Replace(path, "{id}", `([^/]+)`, -1)

	// Ensure the pattern matches the entire path from start to finish.
	regex := regexp.MustCompile("^" + pattern + "$")

	return regex
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Attempt direct route retrieval with method and path (fast path)
	directKey := fmt.Sprintf("%s ^%s$", req.Method, req.URL.Path)
	if routeInfo, ok := r.routes[directKey]; ok {
		routeInfo.handler.ServeHTTP(w, req)
		return
	}

	// Fallback to pattern matching (for parameterized routes)
	for key, routeInfo := range r.routes {
		// Check if the method matches before doing more expensive regex matching
		methodAndPattern := strings.Split(key, " ")
		if req.Method != methodAndPattern[0] {
			continue
		}

		matches := routeInfo.pattern.FindStringSubmatch(req.URL.Path)
		if matches == nil {
			continue
		}
		id := matches[1]
		ctx := context.WithValue(req.Context(), "id", id)
		routeInfo.handler.ServeHTTP(w, req.WithContext(ctx))
		return

	}

	// If no matching route was found, return a 404 error.
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (r *Router) SetupRoutes(middleware ...Middleware) {
	for key, routeInfo := range r.routes {
		handler := routeInfo.handler
		for _, m := range middleware {
			handler = m(handler)
		}
		routeInfo.handler = handler
		r.routes[key] = routeInfo
	}
	http.Handle("/", r)
}
