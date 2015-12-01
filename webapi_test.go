package webapi

import (
	"net/http"
	"testing"
)

func TestNewAPI(t *testing.T) {
	api := NewAPI()

	if api.mux == nil {
		t.Fatal("mux was not initialized in api")
	}
}

func TestHandlers(t *testing.T) {
	code, _ := Handlers(makeRequest(), []Handler{
		func(r *Request) (int, Response) {
			return Next()
		},
		handler,
	})

	if code != 201 {
		t.Fatalf("expected '%d', got '%d'", 201, code)
	}
}

func TestNext(t *testing.T) {
	code, _ := Next()

	if code != 0 {
		t.Fatalf("expected '%d', got '%d'", 0, code)
	}
}

func TestApply(t *testing.T) {
	h := Apply(handler, middleware)
	code, _ := h(makeRequest())

	if code != 418 {
		t.Fatalf("expected '%d', got '%d'", 418, code)
	}
}

func TestWebAPIApply(t *testing.T) {
	api := NewAPI()
	api.Apply(middleware)

	if len(api.middleware) != 1 {
		t.Fatalf("expected API to have %d global middleware, had %d", 1, len(api.middleware))
	}
}

// some helper functions

func handler(r *Request) (int, Response) {
	return 201, Response{}
}

func middleware(handler Handler) Handler {
	return func(r *Request) (int, Response) {
		_, data := handler(r)
		return 418, data
	}
}

func makeRequest() *Request {
	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		panic(err)
	}
	return &Request{*req, nil}
}
