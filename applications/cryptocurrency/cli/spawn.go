package cli

import (
	"errors"
	"fmt"
	"time"

	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	commands_bitcoin "github.com/xmnservices/xmnsuite/applications/cryptocurrency/commands/bitcoin"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/meta"
	webserver "github.com/xmnservices/xmnsuite/applications/cryptocurrency/web"
	"github.com/xmnservices/xmnsuite/configs"
	"github.com/xmnservices/xmnsuite/helpers"
)

func spawn() *cliapp.Command {
	return &cliapp.Command{
		Name:    "spawn",
		Aliases: []string{"s"},
		Usage:   "Spawns a cryptocurrency blockchain",
		Subcommands: []cliapp.Command{
			{
				Name:    "bitcoin",
				Aliases: []string{"f"},
				Usage:   "This is the bitcoin representation blockchain",
				Flags: []cliapp.Flag{
					cliapp.IntFlag{
						Name:  "port",
						Value: 26657,
						Usage: "this is the blockchain port",
					},
					cliapp.IntFlag{
						Name:  "wport",
						Value: 80,
						Usage: "this is the web port",
					},
					cliapp.StringFlag{
						Name:  "dir",
						Value: "./blockchain",
						Usage: "this is the blockchain database path",
					},
					cliapp.StringFlag{
						Name:  "pass",
						Value: "",
						Usage: "this is the password used to decrypt your configuration file",
					},
					cliapp.StringFlag{
						Name:  "file",
						Value: "",
						Usage: "this is the path of your encrypted configuration file",
					},
				},
				Action: func(c *cliapp.Context) error {

					// create the repository:
					repository := configs.SDKFunc.CreateRepository()

					// retrieve the configs:
					retConf, retConfErr := repository.Retrieve(c.String("file"), c.String("pass"))
					if retConfErr != nil {
						return retConfErr
					}

					// spawn the node:
					node := commands_bitcoin.SDKFunc.Spawn(commands_bitcoin.SpawnParams{
						Pass:     c.String("pass"),
						Filename: c.String("file"),
						Dir:      c.String("dir"),
						Port:     c.Int("port"),
					})

					// retrieve the client:
					client, clientErr := node.GetClient()
					if clientErr != nil {
						panic(clientErr)
					}

					// spawn the web server:
					web := webserver.SDKFunc.Create(webserver.CreateParams{
						Port:   c.Int("wport"),
						Client: client,
						Meta:   meta.SDKFunc.CreateMetaData(),
						PK:     retConf.WalletPK(),
					})

					// start the web server:
					err := web.Start()
					if err != nil {
						str := fmt.Sprintf("There was an error while starting the web server: %s", err.Error())
						helpers.Print(str)
					}

					// sleep 1 second before listening to keyboard:
					time.Sleep(time.Second * 1)
					termErr := term.Init()
					if termErr != nil {
						str := fmt.Sprintf("there was an error while enabling the keyboard listening: %s", termErr.Error())
						return errors.New(str)
					}
					defer term.Close()

					// blockchain started, loop until we stop:
					str := fmt.Sprintf("XMN main blockchain spawned, IP: %s\nPress Esc to stop...", client.IP())
					helpers.Print(str)

				keyPressListenerLoop:
					for {
						switch ev := term.PollEvent(); ev.Type {
						case term.EventKey:
							switch ev.Key {
							case term.KeyEsc:
								break keyPressListenerLoop
							}
							break
						}
					}

					// returns:
					return nil
				},
			},
		},
	}
}
