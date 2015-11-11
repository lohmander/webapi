package webapi

import (
	"net/http"
)

type Request struct {
	http.Request
	params map[string]string
}

func WrapHttpRequest(r *http.Request) Request {
	params := make(map[string]string)
	return Request{*r, params}
}

func (r *Request) Param(name string) string {
	return r.params[name]
}
