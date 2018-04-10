package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func hookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	clientIP := getClientIP(r)
	if clientIP != strings.Split(r.RemoteAddr, ":")[0] {
		log.Printf("Received hook for id '%s' from %s on %s\n", id, clientIP, r.RemoteAddr)
	} else {
		log.Printf("Received hook for id '%s' from %s\n", id, r.RemoteAddr)
	}
	rb, err := NewRunBook(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	remoteIP := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
	if !rb.AddrIsAllowed(remoteIP) {
		log.Printf("Hook id '%s' is not allowed from %v\n", id, r.RemoteAddr)
		http.Error(w, "Not authorized.", http.StatusUnauthorized)
		return
	}
	interoplatePOSTData(rb, r)
	interpolateGETData(rb, r)
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

func interpolateGETData(rb *runBook, r *http.Request) {
	vals := r.URL.Query()
	if len(vals) == 0 {
		return
	}
	for k, v := range vals {
		for i := range rb.Scripts {
			for j := range rb.Scripts[i].Args {
				rb.Scripts[i].Args[j] = strings.Replace(rb.Scripts[i].Args[j], "{{"+k+"}}", v[0], -1)
			}
		}
	}
}

func interoplatePOSTData(rb *runBook, r *http.Request) {
	if r.ContentLength == 0 {
		return
	}
	data, err := ioutil.ReadAll(r.Body)
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

func getClientIP(r *http.Request) string {
	remoteIP := strings.Split(r.RemoteAddr, ":")[0]
	if !proxy {
		return remoteIP
	}
	headerVal := r.Header.Get(proxyHeader)
	// proxies can chain upstream client addresses- take only the closest (last) address
	// http://en.wikipedia.org/wiki/X-Forwarded-For
	upstreams := strings.Split(headerVal, ", ")
	return upstreams[len(upstreams)-1]
}
