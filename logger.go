package main

import (
	"log"
	"net/http"
	"time"
)

// We are going to pass our handler to this function,
// which will then wrap the passed handler with logging
// and timing functionality.
func Logger(h func(http.ResponseWriter, *http.Request), n string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			n,
			time.Since(start))

		h(w, r)
	}
}
