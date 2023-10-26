package infra

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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
	path    string
}

type Router struct {
	routes map[string]internalRoute
}

type Middleware func(http.Handler) http.Handler

func NewRouter(dbs map[string]*sql.DB, routes []Route) *Router {
	r := &Router{routes: make(map[string]internalRoute)}
	for _, route := range routes {
		k := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.routes[k] = internalRoute{
			handler: route.Handler(dbs),
			method:  route.Method,
			path:    route.Path,
		}
	}
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	directKey := fmt.Sprintf("%s ^%s$", req.Method, req.URL.Path)
	if routeInfo, ok := r.routes[directKey]; ok {
		routeInfo.handler.ServeHTTP(w, req)
		return
	}
	lookupKey, id := keyAndID(req)
	routeInfo, found := r.routes[lookupKey]
	if !found {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	if id != "" {
		ctx := context.WithValue(req.Context(), "id", id)
		routeInfo.handler.ServeHTTP(w, req.WithContext(ctx))
		return
	}
	routeInfo.handler.ServeHTTP(w, req)
	return

}

func keyAndID(req *http.Request) (string, string) {
	pathSegments := strings.Split(strings.Trim(req.URL.Path, "/"), "/") // Trim is used to remove any leading or trailing slashes
	switch len(pathSegments) {
	case 2: // Potentially /<resource_name>/{id}
		resource := pathSegments[0]
		reconstructedPath := fmt.Sprintf("/%s/{id}", resource)
		return fmt.Sprintf("%s %s", req.Method, reconstructedPath), pathSegments[1]
	case 3: // Potentially /<resource_name>/{id}/<action>
		resource := pathSegments[0]
		action := pathSegments[2]
		reconstructedPath := fmt.Sprintf("/%s/{id}/%s", resource, action)
		return fmt.Sprintf("%s %s", req.Method, reconstructedPath), pathSegments[1]
	default:
		// If the URL doesn't match expected structures, it could be a direct match or a 404.
		return fmt.Sprintf("%s %s", req.Method, req.URL.Path), ""
	}
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
