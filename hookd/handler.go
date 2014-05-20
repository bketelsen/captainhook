package hookd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/bketelsen/captainhook/types"
	"github.com/gorilla/mux"
)

func hookhandler(w http.ResponseWriter, r *http.Request) {
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

func execScript(s types.Script) error {
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

func getScriptFromKey(key string) types.Orchestration {
	p := fmt.Sprintf("%s/%s.json", configdir, key)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Printf("Error opening %s\n", p)
	}
	var o types.Orchestration
	err = json.Unmarshal(b, &o)
	if err != nil {
		log.Println(err)
	}
	return o
}
