package webapi

import (
	"net/http"
	"regexp"
)

// HandlerFunc which is used to respond to a HTTP request.
type handlerFunc func(http.ResponseWriter, *Request)

// Route represents a single route and its handler function.
type Route struct {
	path    string
	handler HandlerFunc
}

// Mux holds all the API routes.
type Mux struct {
	routes []Route
}

// matchPath matches request path against all specified routes.
// If a route is matched, all of the URL parameters will be parsed
// and added to the `Request`-object.
func (mux *Mux) matchPath(path string, route Route) (bool, map[string]string) {
	regex := regexp.MustCompile(route.path)

	if !regex.MatchString(path) {
		return false, nil
	}

	match := regex.FindStringSubmatch(path)
	result := make(map[string]string)

	for i, name := range regex.SubexpNames() {
		if i < 1 {
			continue
		}

		result[name] = match[i]
	}

	return true, result
}

// HandleFunc adds a handler function for given route/path.
func (mux *Mux) HandleFunc(path string, handler handlerFunc) {
	route := Route{path, handler}
	mux.routes = append(mux.routes, route)
}

// ServeHTTP is required to implement the net/http Handler interface
// and thus be compatible with the standard net/http library.
//
// Takes care of serving all HTTP requests.
func (mux *Mux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	for _, route := range mux.routes {
		if match, params := mux.matchPath(path, route); match {
			route.handler(rw, &Request{*r, params})
			return
		}
	}

	http.NotFound(rw, r)
}
