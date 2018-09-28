package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCodeWithWrongHttpMethod_ExpectNotFound(t *testing.T) {
	req, err := http.NewRequest("PUT", "/shorten", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetShortenCode)

	handler.ServeHTTP(rr, req)

	assertWhenError(rr, t)
}

func TestGetCodeNotPresent_ExpectNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/NotExistingCode", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetShortenCode)

	handler.ServeHTTP(rr, req)

	assertWhenError(rr, t)
}

func assertWhenError(rr *httptest.ResponseRecorder, t *testing.T) {
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	if mediatype := rr.Header().Get("Content-Type"); mediatype != "application/json" {
		t.Errorf("handler returned wrong mediatype: got %v want %v",
			mediatype,
			"application/json")
	}
	resp := ApiError{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Error("Handler must respond with valid content")
	}
}
func TestPostNewCodeWithWrongHttpMethod_ExpectNotFound(t *testing.T) {
	var code = "{\"url\":\"www.example.com\",\"shortcode\":\"exp\"}"
	req, err := http.NewRequest("PUT", "/shorten", strings.NewReader(code))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateShortCode)

	handler.ServeHTTP(rr, req)

	assertWhenError(rr, t)
}

func TestPostNewCodeNotPresent_ExpectCreated(t *testing.T) {
	var code = "{\"url\":\"www.example.com\",\"shortcode\":\"exp\"}"
	req, err := http.NewRequest("POST", "/shorten", strings.NewReader(code))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateShortCode)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	if mediatype := rr.Header().Get("Content-Type"); mediatype != "application/json" {
		t.Errorf("handler returned wrong mediatype: got %v want %v",
			mediatype,
			"application/json")
	}
}
