package webapi

import (
	"net/http"
)

type Request struct {
	http.Request
	params map[string]string
}

func (r *Request) Param(name string) string {
	return r.params[name]
}
