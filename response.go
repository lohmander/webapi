package webapi

// Response object that will be marshalled to JSON
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}
