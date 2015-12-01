package webapi

import (
	"encoding/json"
	"net/http"
)

// HTTP methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// Handler function type used to handle HTTP requests.
type Handler func(*Request) (int, Response)

// Middleware is any function that takes and returns a handler.
type Middleware func(Handler) Handler

// GetSupported is the interface that a resource has
// to implement in order to receive GET HTTP requests.
type GetSupported interface {
	Get(*Request) (int, Response)
}

// PostSupported is the interface that a resource has
// to implement in order to receive POST HTTP requests.
type PostSupported interface {
	Post(*Request) (int, Response)
}

// PutSupported is the interface that a resource has
// to implement in order to receive PUT HTTP requests.
type PutSupported interface {
	Put(*Request) (int, Response)
}

// DeleteSupported is the interface that a resource has
// to implement in order to receive DELETE HTTP requests.
type DeleteSupported interface {
	Delete(*Request) (int, Response)
}

// An API manages the resource endpoints that has been added
// to it and routes requests to the appropriate handler function(s).
//
// Since WebApi implements the net/http Handler interface, you can
// add any number of APIs under the same port/server, prefix your
// API and so on. For instance
//
// 	http.Handle("/api/v1", apiv1)
// 	http.Handle("/api/v2", apiv2)
type WebApi struct {
	mux        *Mux
	middleware []Middleware
}

// NewAPI returns a new instance of WebApi.
func NewAPI() *WebApi {
	return &WebApi{&Mux{}, nil}
}

// Handlers lets an enpoint return multiple handlers and thus
// return different responses based on whatever conditional you
// may come up with.
func Handlers(r *Request, handlers []Handler) (int, Response) {
	var code int
	var data Response
	for _, handler := range handlers {
		code, data = handler(r)

		if code > 0 {
			return code, data
		}
	}
	return code, data
}

// Returns a zero status code and empty response. Will make the
// `Handlers` function move on to the next handler.
func Next() (int, Response) {
	return 0, Response{}
}

// Apply applies all the given middleware to provided
// handler function and then returns it.
func Apply(handler Handler, middleware ...Middleware) Handler {
	h := handler
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func (webApi *WebApi) requestHandler(resource interface{}, middleware ...Middleware) HandlerFunc {
	return func(rw http.ResponseWriter, r *Request) {
		var handler Handler

		switch r.Method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				handler = resource.Get
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				handler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				handler = resource.Put
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				handler = resource.Delete
			}
		}

		if handler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		for _, m := range append(webApi.middleware, middleware...) {
			handler = m(handler)
		}

		code, data := handler(r)
		content, err := json.Marshal(data)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

// Apply applies middleware to all subsequently added resources.
func (webApi *WebApi) Apply(middleware ...Middleware) {
	webApi.middleware = middleware
}

// Add adds a resource to the API.
func (webApi *WebApi) Add(path string, resource interface{}, middleware ...Middleware) {
	webApi.mux.HandleFunc(path, webApi.requestHandler(resource, middleware...))
}

// ServeHTTP aliases the webApi.mux.ServeHTTP function.
func (webApi *WebApi) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	webApi.mux.ServeHTTP(rw, r)
}
