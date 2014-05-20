package hookd

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hookhandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	log.Printf("Received Hook for key '%s'\n", key)
	w.Write([]byte("Hello " + key))
}
