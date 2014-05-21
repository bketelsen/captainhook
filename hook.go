package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	log.Printf("Received Hook for key '%s'\n", key)
	rs, err := executeScriptsById(key)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	if echo {
		data, err := json.MarshalIndent(Results{rs}, "", "  ")
		if err != nil {
			log.Println(err.Error())
		}
		w.Write(data)
	}
}
