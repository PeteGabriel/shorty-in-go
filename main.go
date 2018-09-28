package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/shortid"
)

var mediatype = "application/json"

func main() {

	a := App{}
	//TODO this should go into env vars
	a.Initialize("dummy0", "dummy1", "dummy2")
	a.Run(":8080")
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
