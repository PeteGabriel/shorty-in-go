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

	ResNotFound    = "Resource not found"
	ContentHeader  = "Content-Type"
	MediaType      = "application/json"
	LocationHeader = "Location"

	ShouldBeValidCode = "The shortcode provided must comply with ([0-9a-zA-Z]+)$"
	CodeNotFound      = "The shortcode cannot be found in the system"
	CodeIsInvalid     = "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{4,}$."
)

var codes map[string]Code

type App struct {
	Router *mux.Router
}

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
		writeCustomHeader(w, http.StatusNotFound)
		resp := buildErrMsg(http.StatusNotFound, ResNotFound)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}
	// map of parameters
	sc := mux.Vars(r)["shortcode"]

	vreg := regexp.MustCompile("([0-9a-zA-Z]+)$")
	if m := vreg.FindStringSubmatch(sc); m != nil {
		if v, ok := codes[m[1]]; ok {
			w.Header().Set(ContentHeader, MediaType)
			w.Header().Set(LocationHeader, v.Url)
			w.WriteHeader(http.StatusFound)
			return
		}
		//TODO add better response
		writeCustomHeader(w, http.StatusNotFound)
		resp := buildErrMsg(http.StatusNotFound, CodeNotFound)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}
	//BadRequest since code provided does not comply with regex
	writeCustomHeader(w, http.StatusBadRequest)
	resp := buildErrMsg(http.StatusBadRequest, ShouldBeValidCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
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
		resp := buildErrMsg(http.StatusUnprocessableEntity, CodeIsInvalid)

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
	w.Header().Set(ContentHeader, MediaType)
	w.WriteHeader(code)
}

func handleInvalidMethod(w http.ResponseWriter) {
	writeCustomHeader(w, http.StatusNotFound)
	resp := buildErrMsg(http.StatusNotFound, ResNotFound)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func buildErrMsg(code int, desc string) APIError {
	return APIError{
		Error: code,
		Desc:  desc,
	}
}
