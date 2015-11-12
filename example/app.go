package main

import (
	"fmt"
	"net/http"

	"github.com/lohmander/webapi"
)

func main() {
	api := webapi.NewAPI()
	api.Apply(Logger)
	api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot)

	http.Handle("/api/", api)
	http.ListenAndServe(":3002", nil)
}

type Item struct{}

func (item Item) Post(request *webapi.Request) (int, webapi.Response) {
	var body interface{}

	err := request.UnmarshalBody(&body)
	if err != nil {
		return 500, webapi.Response{
			Error: err,
		}
	}

	return 200, webapi.Response{
		Data: map[string]interface{}{
			"body":    body,
			"idParam": request.Param("id"),
		},
	}
}

// some middleware

func Logger(handler webapi.Handler) webapi.Handler {
	return func(r *webapi.Request) (int, webapi.Response) {
		code, data := handler(r)
		fmt.Println(code, r.Method, r.URL.Path)
		return code, data
	}
}

func Teapot(handler webapi.Handler) webapi.Handler {
	return func(r *webapi.Request) (int, webapi.Response) {
		_, data := handler(r)
		return 418, data
	}
}
