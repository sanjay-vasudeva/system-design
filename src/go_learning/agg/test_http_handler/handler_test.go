package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	handler(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d", http.StatusOK, rr.Result().StatusCode)
	}
	defer rr.Result().Body.Close()

	expected := "FOO"
	actual := rr.Body.String()
	//
	if actual != expected {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}
