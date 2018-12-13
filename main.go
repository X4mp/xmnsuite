package main

import (
	"log"
	"os"

	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	forex "github.com/xmnservices/xmnsuite/applications/forex"
)

func reset() {
	term.Sync()
}

func main() {
	app := cliapp.NewApp()
	app.Version = "2018.12.13"
	app.Name = "xmn"
	app.Usage = "The XMN network is a network of blockchain applications used to decentrealized businesses"
	app.Commands = []cliapp.Command{
		{
			Name:        "forex",
			Aliases:     []string{"f"},
			Usage:       "This is the forex blockchain application",
			Subcommands: forex.SDKFunc.Create(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
