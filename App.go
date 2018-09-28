package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     map[string]Code
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
