package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.Printf("Received hook for id '%s' from %s\n", id, r.RemoteAddr)
	rb, err := NewRunBook(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	response, err := rb.execute()
	if err != nil {
        log.Println(err.Error())
        http.Error(w, err.Error(), 500)
        return
    }
	if echo {
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Println(err.Error())
		}
		w.Write(data)
	}
}
