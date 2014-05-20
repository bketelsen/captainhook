package createconfig

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/bketelsen/captainhook/types"
	"github.com/robmerrell/comandante"
)

var filename string

func NewCommand() *comandante.Command {
	return comandante.NewCommand("createconfig", "create a command configuration template", func() error {

		return createCommand()
	})
}

func createCommand() error {

	o := types.Orchestration{}

	s1 := types.Script{Command: "ls"}
	s2 := types.Script{Command: "ps"}

	scripts := []types.Script{s1, s2}
	o.Scripts = scripts

	fmt.Printf("Some Config would be spit out here and it would be named %s\n", filename)

	output, _ := json.MarshalIndent(o, "", "    ")
	fmt.Println(string(output))
	return nil
}

func GetFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&filename, "filename", "sample.json", "File to write")
}
