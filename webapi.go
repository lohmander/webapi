package webapi

import (
	"encoding/json"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Handler func(*Request) (int, Response)
type Middleware func(Handler) Handler

type GetSupported interface {
	Get(*Request) (int, Response)
}

type PostSupported interface {
	Post(*Request) (int, Response)
}

type PutSupported interface {
	Put(*Request) (int, Response)
}

type DeleteSupported interface {
	Delete(*Request) (int, Response)
}

type WebApi struct {
	mux        *Mux
	middleware []Middleware
}

func NewAPI() *WebApi {
	return &WebApi{&Mux{}, nil}
}

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

func Next() (int, Response) {
	return 0, Response{}
}

func Bake(handler Handler, middleware ...Middleware) Handler {
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

func (webApi *WebApi) Apply(middleware ...Middleware) {
	webApi.middleware = middleware
}

func (webApi *WebApi) Add(path string, resource interface{}, middleware ...Middleware) {
	webApi.mux.HandleFunc(path, webApi.requestHandler(resource, middleware...))
}

func (webApi *WebApi) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	webApi.mux.ServeHTTP(rw, r)
}
