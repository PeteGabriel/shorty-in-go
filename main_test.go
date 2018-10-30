package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCode_WithWrongHttpMethod_ExpectNotFound(t *testing.T) {
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

func TestGetCodeNotCompliant_ExpectBadRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/_", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetShortenCode)

	handler.ServeHTTP(rr, req)

	assertWhenError(rr, t)
}

func TestPostNewCode_WithWrongHttpMethod_ExpectNotFound(t *testing.T) {
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
func TestPostNewCode_NotPresent_ExpectCreated(t *testing.T) {
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

	checkMediatype(t, rr.Header().Get("Content-Type"))
}
func TestPostNewCode_BadRequest_ExpectBadRequest(t *testing.T) {

}

/*********************************
* Utility methods
**********************************/

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkMediatype(t *testing.T, actual string) {
	if actual != "application/json" {
		t.Errorf("handler returned wrong mediatype: got %v want %v",
			actual, "application/json")
	}
}

//TODO write a wrapper to encapsulate code usage
// assertWhenError(rr, t, http.StatusNotFound)

func assertWhenError(rr *httptest.ResponseRecorder, t *testing.T) {
	checkResponseCode(t, http.StatusNotFound, rr.Code)

	checkMediatype(t, rr.Header().Get("Content-Type"))

	resp := APIError{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Error("Handler must respond with valid content")
	}

	if resp.Desc != "Resource not found" {
		t.Error("Response description field its not correct.")
	}
}
