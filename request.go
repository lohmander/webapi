package webapi

import (
	"encoding/json"
	"net/http"
)

// Request object which inherits from net/http Request.
// Also holds URL params that has been parsed according to the
// route regex which can be accesed through the Param function.
type Request struct {
	http.Request
	params map[string]string
}

// Param returns the URL param value for given key/name.
func (r *Request) Param(name string) string {
	return r.params[name]
}

// UnmarshalBody unmarshals request body from JSON to the provided target.
func (r *Request) UnmarshalBody(target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}
