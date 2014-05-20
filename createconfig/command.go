package createconfig

import (
	"fmt"

	"github.com/robmerrell/comandante"
)

func NewCommand() *comandante.Command {
	return comandante.NewCommand("createconfig", "create a command configuration template", func() error {

		return createCommand()
	})
}

func createCommand() error {
	fmt.Println("Some Config would be spit out here")
	return nil
}
