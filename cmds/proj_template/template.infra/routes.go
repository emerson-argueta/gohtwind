package infra

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Route struct {
	Method     string
	HandleFunc http.HandlerFunc
	Path       string
	handler    http.Handler
}

var baseRouter *Router

type Router struct {
	http.Handler
	routes map[string]Route
}

type Middleware func(http.Handler) http.Handler

func init() {
	baseRouter = &Router{routes: make(map[string]Route)}
}
func NewRouter(routes []Route) *Router {
	r := &Router{routes: make(map[string]Route)}
	for _, route := range routes {
		k := fmt.Sprintf("%s %s", route.Method, route.Path)
		route.handler = route.HandleFunc
		r.routes[k] = route
		baseRouter.routes[k] = route
	}
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodOverride := req.FormValue("_method")
	directKey := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	switch methodOverride {
	case "PATCH":
		directKey = fmt.Sprintf("%s %s", "PATCH", req.URL.Path)
	case "DELETE":
		directKey = fmt.Sprintf("%s %s", "DELETE", req.URL.Path)
	}
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
	method := req.Method
	methodOverride := req.FormValue("_method")
	switch methodOverride {
	case "PATCH":
		method = "PATCH"
	case "DELETE":
		method = "DELETE"
	}
	pathSegments := strings.Split(strings.Trim(req.URL.Path, "/"), "/") // Trim is used to remove any leading or trailing slashes
	switch len(pathSegments) {
	case 2: // Potentially /<feature_name>/{id}
		resource := pathSegments[0]
		reconstructedPath := fmt.Sprintf("/%s/{id}", resource)
		return fmt.Sprintf("%s %s", method, reconstructedPath), pathSegments[1]
	case 3: // Potentially /<feature_name>/{id}/<action> or /<feature_name>/<resource_name>/{id}
		resource := pathSegments[0]
		actionOrId := pathSegments[2]
		// id is a number
		if _, err := strconv.Atoi(actionOrId); err == nil {
			reconstructedPath := fmt.Sprintf("/%s/%s/{id}", pathSegments[0], pathSegments[1])
			return fmt.Sprintf("%s %s", method, reconstructedPath), actionOrId
		}
		reconstructedPath := fmt.Sprintf("/%s/{id}/%s", resource, actionOrId)
		return fmt.Sprintf("%s %s", method, reconstructedPath), pathSegments[1]
	case 4: // Potentially /<feature_name>/<resource_name>/{id}/<action>
		reconstructedPath := fmt.Sprintf("/%s/%s/{id}/%s", pathSegments[0], pathSegments[1], pathSegments[3])
		return fmt.Sprintf("%s %s", method, reconstructedPath), pathSegments[2]
	default:
		// If the URL doesn't match expected structures, it could be a direct match or a 404.
		return fmt.Sprintf("%s %s", method, req.URL.Path), ""
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
		baseRouter.routes[key] = routeInfo
	}
}

func ActivateRoutes() {
	http.Handle("/", baseRouter)
}
