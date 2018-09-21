package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
)

var mediatype = "application/json; charset=UTF-8"

func main() {
	http.HandleFunc("/", Logger(getShortenCode, "Get shortcode by name"))
	http.HandleFunc("/shorten", Logger(newShortCode, "Create new shortcode"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// POST /shorten
func newShortCode(w http.ResponseWriter, r *http.Request) {
	if err := require(r.Method == http.MethodPost, func() { http.NotFound(w, r) }); err != nil {
		return
	}

	vreg := regexp.MustCompile("^/([0-9a-zA-Z_]){4,}$")
	if m := vreg.FindStringSubmatch(r.URL.Path); m == nil {
		writeCustomHeader(w, http.StatusUnprocessableEntity)
		resp := ApiError{
			Error: 402,
			Desc:  "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{4,}$.",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	writeCustomHeader(w, http.StatusCreated)
	//TODO apply right logic
	code := Code{Url: "url", Shortcode: "Code"}
	if err := json.NewEncoder(w).Encode(code); err != nil {
		panic(err)
	}
}

func getShortenCode(w http.ResponseWriter, r *http.Request) {
	if err := require(r.Method == http.MethodGet, func() { http.NotFound(w, r) }); err != nil {
		return
	}

}

// If pred is not true, resolve function ifFalse
// and an error is returned if so.
func require(pred bool, ifFalse func()) error {
	if !pred {
		ifFalse()
		return errors.New("Require clause not valid")
	}
	return nil
}

func writeCustomHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", mediatype)
	w.WriteHeader(code)
}
