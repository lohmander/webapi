package main

import (
	"fmt"
	"net/http"

	"github.com/lohmander/webapi"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello there")
}

func main() {
	api := webapi.NewAPI()
	api.Add("/subscriptions", &Subscription{}, TeapotMiddleware)

	http.ListenAndServe(":3002", api)
}

func TeapotMiddleware(handler webapi.Handler) webapi.Handler {
	return func(r *webapi.Request) (int, webapi.Response) {
		_, data := handler(r)
		return 418, data
	}
}

type Subscription struct{}

func (s Subscription) Post(request *webapi.Request) (int, webapi.Response) {
	var data interface{} = map[string]string{
		"test": "topp",
	}

	return 200, webapi.Response{
		Data: &data,
	}
}
