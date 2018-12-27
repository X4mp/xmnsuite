package main

import (
	"log"
	"os"

	term "github.com/nsf/termbox-go"
	amino "github.com/tendermint/go-amino"
	cliapp "github.com/urfave/cli"
	forex "github.com/xmnservices/xmnsuite/applications/forex"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
)

func reset() {
	term.Sync()
}

func main() {

	// register amino:
	cdc := amino.NewCodec()
	forex.Register(cdc)

	// get the core commands:
	coreCmds := core.SDKFunc.CreateCommands()

	// merge the core to the forex cmds:
	forexCmds := forex.SDKFunc.Create()
	for _, oneCoreCmd := range coreCmds {
		forexCmds = append(forexCmds, oneCoreCmd)
	}

	app := cliapp.NewApp()
	app.Version = "2018.12.13"
	app.Name = "xmn"
	app.Usage = "The XMN network is a network of blockchain applications used to decentrealized businesses"
	app.Commands = []cliapp.Command{
		{
			Name:        "forex",
			Aliases:     []string{"f"},
			Usage:       "This is the forex blockchain application",
			Subcommands: forexCmds,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
