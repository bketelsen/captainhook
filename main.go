package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var configdir string

func init() {
	flag.StringVar(&configdir, "configdir", "", "config dir to use")
}



func main() {
	flag.Parse()
	if configdir == "" {
		log.Fatal("configdir is required")
	}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", hookHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Listening on port %s\n", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
