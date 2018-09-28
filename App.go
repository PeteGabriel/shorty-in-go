package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

// Initialize func should be used to customize
// the server with some options.
func (a *App) Initialize(user, password, dbname string) {
	http.HandleFunc("/shorten", Logger(CreateShortCode, "Create new shortcode"))
	http.HandleFunc("/", Logger(GetShortenCode, "Get shortcode by name"))
}

// Run func should be used to start the server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
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

	if saved, err := SaveCode(code); saved {
		writeCustomHeader(w, http.StatusCreated)
		if err := json.NewEncoder(w).Encode(CodeDto{Shortcode: code.Shortcode}); err != nil {
			panic(err)
		}
	} else {
		//TODO this should be verified if it is actually a "bad request" or something else
		writeCustomHeader(w, http.StatusBadRequest)
	}
}
