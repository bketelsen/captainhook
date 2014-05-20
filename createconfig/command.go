package createconfig

import (
	"flag"
	"fmt"

	"github.com/robmerrell/comandante"
)

var filename string

func NewCommand() *comandante.Command {
	return comandante.NewCommand("createconfig", "create a command configuration template", func() error {

		return createCommand()
	})
}

func createCommand() error {
	fmt.Printf("Some Config would be spit out here and it would be named %s\n", filename)
	return nil
}

func GetFlagHandler(fs *flag.FlagSet) {
	fs.StringVar(&filename, "filename", "sample.json", "File to write")
}
