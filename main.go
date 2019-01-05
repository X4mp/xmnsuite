package main

import (
	"log"
	"os"

	term "github.com/nsf/termbox-go"
	amino "github.com/tendermint/go-amino"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

func reset() {
	term.Sync()
}

func main() {

	// register amino:
	cdc := amino.NewCodec()
	core.Register(cdc)

	// create the meta to generate the request registry:
	meta.SDKFunc.Create(meta.CreateParams{})

	app := cliapp.NewApp()
	app.Version = "2019.01.01"
	app.Name = "xmn"
	app.Usage = "This is the xmn core application"
	app.Commands = []cliapp.Command{
		*cli.SDKFunc.Spawn(),
		*cli.SDKFunc.Wallet(),
		*cli.SDKFunc.Genesis(),
		*cli.SDKFunc.Information(),
		*cli.SDKFunc.Balance(),
		*cli.SDKFunc.User(),
		*cli.SDKFunc.Transfer(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
