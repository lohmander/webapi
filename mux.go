package webapi

import (
	"net/http"
	"regexp"
)

type HandlerFunc func(http.ResponseWriter, *Request)

type Route struct {
	path    string
	handler HandlerFunc
}

type Mux struct {
	routes []Route
}

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

func (mux *Mux) HandleFunc(path string, handler HandlerFunc) {
	route := Route{path, handler}
	mux.routes = append(mux.routes, route)
}

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
