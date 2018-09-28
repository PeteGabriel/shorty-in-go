package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/shortid"
)

var mediatype = "application/json"

var codes map[string]Code

func init() {
	codes = make(map[string]Code)
}

func main() {

	a := App{}
	//TODO this should go into env vars
	a.Initialize("dummy0", "dummy1", "dummy2")
	a.Run(":8080")
}

// POST /shorten
func CreateShortCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleInvalidMethod(w)
		return
	}

	vreg := regexp.MustCompile("^/([0-9a-zA-Z_]){4,}$")
	if m := vreg.FindStringSubmatch(r.URL.Path); m == nil {
		writeCustomHeader(w, http.StatusUnprocessableEntity)
		resp := APIError{
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
func GetShortenCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		//TODO change to return api error
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
		resp := APIError{
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

func handleInvalidMethod(w http.ResponseWriter) {
	writeCustomHeader(w, http.StatusNotFound)
	resp := APIError{
		Error: http.StatusNotFound,
		Desc:  "Resource not found",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
