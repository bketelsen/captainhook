package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	configdir   string
	echo        bool
	proxy       bool
	version     bool
	proxyHeader string
	listenAddr  string
)

func init() {
	flag.StringVar(&configdir, "configdir", "", "config dir to use")
	flag.BoolVar(&echo, "echo", false, "send output from script")
	flag.BoolVar(&proxy, "enable-proxy", false, "enable parsing of X-Forwarded-For header")
	flag.BoolVar(&version, "version", false, "display version")
	flag.StringVar(&proxyHeader, "proxy-header", "X-Forwarded-For", "header to use for upstream client IP")
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:8080", "http listen address")
}

func main() {
	flag.Parse()

	if BuildDate == "" {
		BuildDate = time.Now().Format("15:04:05 Jan 2, 2006")
	}
	if version {
		fmt.Printf("captain hook version %s built %v\n", Version, BuildDate)

		return
	}

	if configdir == "" {
		log.Fatal("configdir is required")
	}

	r := mux.NewRouter()
	r.HandleFunc("/{id}", hookHandler).Methods("POST")
	http.Handle("/", r)

	log.Printf("Listening on %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
