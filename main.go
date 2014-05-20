package createconfig

import (
	"github.com/bketelsen/captainhook/createconfig"
	"github.com/bketelsen/captainhook/hookd"
	"github.com/bketelsen/captainhook/log"
	"github.com/robmerrell/comandante"
)

func main() {
	captain := comandante.New("captainhook", "")
	captain.IncludeHelp()

	hookd := hookd.NewCommand()
	captain.RegisterCommand(hookd)

	config := createconfig.NewCommand()
	config.FlagInit = createconfig.GetFlagHandler
	captain.RegisterCommand(config)

	log.PanicIf(captain.Run())
}
