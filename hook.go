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

type response struct {
	Results []result `json:"results"`
}

type result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	log.Printf("Received Hook for key '%s'\n", key)
	results := processHandler(key)

	data, err := json.MarshalIndent(&response{results}, "", "  ")
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(data)
}

func processHandler(key string) ([]result) {
	results := make([]result, 0)
	script := getScriptFromKey(key)
	for _, x := range script.Scripts {
		r, err := execScript(x)
		if err != nil {
			log.Println("ERROR :" + err.Error())
		}
		results = append(results, r) 
	}
	return results
}

func execScript(s script) (result, error) {
	cmd := exec.Command(s.Command, s.Args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	r := result{
		stdout.String(),
		stderr.String(),
	}
	return r, err
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
