package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

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
	interoplatePOSTData(rb, r)
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

func interoplatePOSTData(rb *runBook, r *http.Request) {
	if r.ContentLength == 0 {
		return
	}
	data := make([]byte, r.ContentLength)
	_, err := r.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()
	stringData := string(data[:r.ContentLength])
	for i := range rb.Scripts {
		for j := range rb.Scripts[i].Args {
			rb.Scripts[i].Args[j] = strings.Replace(rb.Scripts[i].Args[j], "{{POST}}", stringData, -1)
		}
	}
}
