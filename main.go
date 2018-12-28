package main

import (
	"log"
	"os"

	term "github.com/nsf/termbox-go"
	amino "github.com/tendermint/go-amino"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
)

func reset() {
	term.Sync()
}

func main() {

	// register amino:
	cdc := amino.NewCodec()
	cryptocurrency.Register(cdc)

	// get the core commands:
	coreCmds := core.SDKFunc.CreateCommands()

	// merge the core to the cryptocurrency cmds:
	cryptoCmds := cryptocurrency.SDKFunc.Create()
	for _, oneCoreCmd := range coreCmds {
		cryptoCmds = append(cryptoCmds, oneCoreCmd)
	}

	app := cliapp.NewApp()
	app.Version = "2018.12.28"
	app.Name = "xmn"
	app.Usage = "The XMN network is a network of blockchain applications used to decentralize marketplaces"
	app.Commands = []cliapp.Command{
		{
			Name:        "cryptocurrency",
			Aliases:     []string{"f"},
			Usage:       "This is the cryptocurrency representation blockchain applications",
			Subcommands: cryptoCmds,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
