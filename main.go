package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


var configdir = ""

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", hookHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Listening on port %s\n", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
