package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"syscall"
)

// runBook represents a collection of scripts.
type runBook struct {
	Scripts         []script `json:"scripts"`
	AllowedNetworks Networks `json:"allowedNetworks,omitempty"`
}

type runBookResponse struct {
	Results []result `json:"results"`
}

type result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	StatusCode int    `json:"status_code"`
}

type script struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// Networks is its own struct for JSON unmarshalling gymnastics
type Networks struct {
	Networks []net.IPNet
}

// UnmarshalJSON for custom type Networks
func (nets *Networks) UnmarshalJSON(data []byte) error {
	ns := []string{}
	if err := json.Unmarshal(data, &ns); err != nil {
		return err
	}

	nets.Networks = make([]net.IPNet, len(ns))
	for i, nw := range ns {
		_, ipnet, err := net.ParseCIDR(nw)
		if err != nil {
			return err
		}
		nets.Networks[i] = *ipnet
	}
	return nil
}

// NewRunBook returns the runBook identified by id.
func NewRunBook(id string) (*runBook, error) {
	return getRunBookById(id)
}

func (r *runBook) AddrIsAllowed(remoteIP net.IP) bool {
	if len(r.AllowedNetworks.Networks) == 0 {
		return true
	}
	for _, nw := range r.AllowedNetworks.Networks {
		if nw.Contains(remoteIP) {
			return true
		}
	}
	return false
}

func (r *runBook) execute() (*runBookResponse, error) {
	results := make([]result, 0)
	for _, x := range r.Scripts {
		r, err := execScript(x)
		if err != nil {
			log.Println("ERROR :" + err.Error())
		}
		results = append(results, r)
	}
	return &runBookResponse{results}, nil
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
		cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(),
	}
	return r, err
}

func getRunBookById(id string) (*runBook, error) {
	var r = new(runBook)
	runBookPath := fmt.Sprintf("%s/%s.json", configdir, id)
	data, err := ioutil.ReadFile(runBookPath)
	if err != nil {
		return r, fmt.Errorf("cannot read run book %s: %s", runBookPath, err)
	}
	err = json.Unmarshal(data, r)
	if err != nil {
		return r, err
	}
	return r, nil
}
