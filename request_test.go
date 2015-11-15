package webapi

import (
	"bytes"
	"net/http"
	"testing"
)

func TestParamFunc(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		panic(err)
	}

	params := map[string]string{"test": "testing"}
	request := Request{*req, params}

	if request.Param("test") != "testing" {
		t.Fatalf("expected '%s', got '%s'", "testing", request.Param("test"))
	}
}

func TestUnmarshalBodyFunc(t *testing.T) {
	body := []byte(`{"test": "testing"}`)
	req, err := http.NewRequest("POST", "http://localhost", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	request := Request{*req, nil}
	target := make(map[string]string)
	err = request.UnmarshalBody(&target)
	if err != nil {
		t.Fatal(err)
	}

	if target["test"] != "testing" {
		t.Fatalf("expected '%s', got '%s'", "testing", target["test"])
	}
}
