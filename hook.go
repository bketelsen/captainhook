package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

type orchestration struct {
	Scripts []script `json:"scripts"`
}

type script struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	log.Printf("Received Hook for key '%s'\n", key)
	processHandler(key)
	w.Write([]byte("Hello " + key + " " + configdir))
}

func processHandler(key string) {
	script := getScriptFromKey(key)
	for _, x := range script.Scripts {
		if err := execScript(x); err != nil {
			log.Println("ERROR :" + err.Error())
		}
	}
}

func execScript(s script) error {
	cmd := exec.Command(s.Command, s.Args...)
	var out bytes.Buffer
	// keeping stdout just to see what happened
	// maybe don't need it?  Probably should log it
	//TBD
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(out.String())
	return err
}

func getScriptFromKey(key string) orchestration {
	p := fmt.Sprintf("%s/%s.json", configdir, key)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Printf("Error opening %s\n", p)
	}
	var o orchestration
	err = json.Unmarshal(b, &o)
	if err != nil {
		log.Println(err)
	}
	return o
}
