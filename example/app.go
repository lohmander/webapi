package main

import (
	"fmt"
	"net/http"

	"github.com/lohmander/webapi"
)

func main() {
	api := webapi.NewAPI()
	api.Apply(Logger)
	api.Add(`/subscriptions/(?P<id>\d+)$`, &Subscription{}, Teapot)

	http.Handle("/api/", api)
	http.ListenAndServe(":3002", nil)
}

type Subscription struct{}

func (s Subscription) Post(request *webapi.Request) (int, webapi.Response) {
	var data interface{} = map[string]string{
		"param": request.Param("id"),
	}

	return 200, webapi.Response{
		Data: &data,
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
