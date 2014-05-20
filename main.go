package main

import (
	"github.com/bketelsen/captainhook/createconfig"
	"github.com/bketelsen/captainhook/hookd"
	"github.com/bketelsen/captainhook/log"
	"github.com/robmerrell/comandante"
)

func main() {
	captain := comandante.New("captainhook", "")
	captain.IncludeHelp()

	registerCommand := func(c *comandante.Command) {
		captain.RegisterCommand(c)
	}

	registerCommand(httpd.NewCommand())
	registerCommand(createconfig.NewCommand())

	log.PanicIf(captain.Run())
}
