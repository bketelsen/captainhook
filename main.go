package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	configdir  string
	echo       bool
	listenAddr string
)

func init() {
	flag.StringVar(&configdir, "configdir", "", "config dir to use")
	flag.BoolVar(&echo, "echo", false, "send output from script")
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:8080", "http listen address")
}

func main() {
	flag.Parse()
	if configdir == "" {
		log.Fatal("configdir is required")
	}

	r := mux.NewRouter()
	r.HandleFunc("/{key}", hookHandler).Methods("POST")
	http.Handle("/", r)

	log.Printf("Listening on %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
