package types

import "net/http"

type RouteType []RouteFields

type Middlewares []func(http.Handler) http.Handler

type RouteFields struct {
	Method      string
	Path        string
	Handler     any
	Description string
	Middlewares
}

// Helper function to add a prefix to all routes
func WithPrefix(prefix string, routes RouteType) RouteType {
	for i := range routes {
		routes[i].Path = prefix + routes[i].Path
	}
	return routes
}
