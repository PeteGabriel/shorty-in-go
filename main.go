package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/shortid"
)

var mediatype = "application/json"

var codes map[string]Code

func main() {
	codes = make(map[string]Code)
	http.HandleFunc("/", Logger(getShortenCode, "Get shortcode by name"))
	http.HandleFunc("/shorten", Logger(newShortCode, "Create new shortcode"))

	log.Fatal(http.ListenAndServe(":8000", nil))
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
			Error: http.StatusUnprocessableEntity,
			Desc:  "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{4,}$.",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	data := Code{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		writeCustomHeader(w, http.StatusNotFound)
		return
	}

	code := Code{}
	if data.Url == "" {
		writeCustomHeader(w, http.StatusBadRequest)
		return
	}
	//todo find another way
	if data.Shortcode == "" {
		code.Shortcode = genCode()
	} else {
		code.Shortcode = data.Shortcode
	}
	code.Url = data.Url
	codes[code.Shortcode] = code

	writeCustomHeader(w, http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CodeDto{Shortcode: code.Shortcode}); err != nil {
		panic(err)
	}
}

func getShortenCode(w http.ResponseWriter, r *http.Request) {
	if err := require(r.Method == http.MethodGet, func() { http.NotFound(w, r) }); err != nil {
		return
	}

	vreg := regexp.MustCompile("^/([0-9a-zA-Z_]+)$")
	if m := vreg.FindStringSubmatch(r.URL.Path); m != nil {
		if v, ok := codes[m[1]]; ok {
			http.Redirect(w, r, v.Url, http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	} else {
		http.NotFound(w, r)
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

func genCode() string {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		panic(err)
	}

	c, err := sid.Generate()
	if err != nil {
		panic(err)
	}

	return c
}
