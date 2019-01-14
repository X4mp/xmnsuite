package main

import (
	"log"
	"net/http"
	"os"

	update "github.com/inconshreveable/go-update"
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

	// create the meta:
	meta.SDKFunc.Create(meta.CreateParams{})

	app := cliapp.NewApp()
	app.Version = "2019.01.01"
	app.Name = "xmn"
	app.Usage = "This is the xmn core application"
	app.Commands = []cliapp.Command{
		*cli.SDKFunc.Spawn(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	upErr := update.Apply(resp.Body, update.Options{})
	return upErr
}
