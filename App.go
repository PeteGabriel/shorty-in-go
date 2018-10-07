package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/shortid"
)

const (
	GET  = "GET"
	POST = "POST"
)

var codes map[string]Code

type App struct {
	Router *mux.Router
}

var mediatype = "application/json"

func init() {
	codes = make(map[string]Code)
}

// Initialize func should be used to customize
// the server with some options.
func (a *App) Initialize(user, password, dbname string) {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/shorten", CreateShortCode).Methods(POST)
	a.Router.HandleFunc("/{shortcode}", GetShortenCode).Methods(GET)
}

// Run func should be used to start the server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//GET /:code
func GetShortenCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		//TODO change to return api error
		http.NotFound(w, r)
		return
	}
	// map of parameters
	sc := mux.Vars(r)["shortcode"]

	vreg := regexp.MustCompile("([0-9a-zA-Z_]+)$")
	if m := vreg.FindStringSubmatch(sc); m != nil {
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
	//TODO Handle the fact that reqeust was poorly made
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
		writeCustomHeader(w, http.StatusBadRequest)
		return
	}

	if code.Url == "" {
		//TODO this should be a specific error, instead of soemthing so generic as error type
		writeCustomHeader(w, http.StatusBadRequest)
		return
	}

	if code.Shortcode == "" {
		code.Shortcode = genCode()
	}

	codes[code.Shortcode] = code
	writeCustomHeader(w, http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CodeDto{Shortcode: code.Shortcode}); err != nil {
		panic(err)
	}

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

func writeCustomHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", mediatype)
	w.WriteHeader(code)
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
