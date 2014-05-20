package hookd

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robmerrell/comandante"
)

func NewCommand() *comandante.Command {
	return comandante.NewCommand("httpd", "run http server", func() error {
		r := mux.NewRouter()
		r.HandleFunc("/{key}", hookhandler).Methods("GET")
		http.Handle("/", r)
		log.Printf("Listening on port %s\n", ":8080")
		return http.ListenAndServe(":8080", nil)
	})
}
