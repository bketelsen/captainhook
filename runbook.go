package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"
)

type Results struct {
	Results []Result `json:"results"`
}

type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	StatusCode int    `json:"status_code"`
}

type runBook struct {
	Scripts []script `json:"scripts"`
}

type script struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func executeScriptsById(key string) ([]Result, error) {
	rb, err := getRunBookFromKey(key)
	if err != nil {
		return nil, err
	}
	rs := make([]Result, 0)
	for _, x := range rb.Scripts {
		r, err := execScript(x)
		if err != nil {
			log.Println("ERROR :" + err.Error())
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func execScript(s script) (Result, error) {
	cmd := exec.Command(s.Command, s.Args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	r := Result{
		stdout.String(),
		stderr.String(),
		cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(),
	}
	return r, err
}

func getRunBookFromKey(key string) (runBook, error) {
	var r runBook
	runBookPath := fmt.Sprintf("%s/%s.json", configdir, key)
	data, err := ioutil.ReadFile(runBookPath)
	if err != nil {
		return r, fmt.Errorf("cannot read run book %s: %s", runBookPath, err)
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}
