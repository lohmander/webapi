package webapi

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	http.Request
	params map[string]string
}

func (r *Request) Param(name string) string {
	return r.params[name]
}

func (r *Request) UnmarshalBody(target *interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}
