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

	hook := hookd.NewCommand()
	hook.FlagInit = hookd.GetFlagHandler
	captain.RegisterCommand(hook)

	config := createconfig.NewCommand()
	config.FlagInit = createconfig.GetFlagHandler
	captain.RegisterCommand(config)

	log.PanicIf(captain.Run())
}
