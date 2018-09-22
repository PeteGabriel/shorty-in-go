package main

import (
	"encoding/json"
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
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
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

	code := Code{}
	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		writeCustomHeader(w, http.StatusNotFound)
		return
	}

	if code.Url == "" {
		writeCustomHeader(w, http.StatusBadRequest)
		return
	}
	//todo find another way
	if code.Shortcode == "" {
		code.Shortcode = genCode()
	}
	codes[code.Shortcode] = code

	writeCustomHeader(w, http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CodeDto{Shortcode: code.Shortcode}); err != nil {
		panic(err)
	}
}

//GET /:code
func getShortenCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	vreg := regexp.MustCompile("^/([0-9a-zA-Z_]+)$")
	if m := vreg.FindStringSubmatch(r.URL.Path); m != nil {
		if v, ok := codes[m[1]]; ok {
			w.Header().Set("Content-Type", mediatype)
			w.Header().Set("Location", v.Url)
			w.WriteHeader(http.StatusFound)
			return
		}
		//TODO add better response
		writeCustomHeader(w, http.StatusNotFound)
		resp := ApiError{
			Error: http.StatusNotFound,
			Desc:  "The shortcode cannot be found in the system",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	http.NotFound(w, r)
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
