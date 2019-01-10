package cli

import (
	"errors"
	"fmt"
	"time"

	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/commands"
	"github.com/xmnservices/xmnsuite/helpers"
)

func spawn() *cliapp.Command {
	return &cliapp.Command{
		Name:    "spawn",
		Aliases: []string{"s"},
		Usage:   "Spawns the core blockchain",
		Flags: []cliapp.Flag{
			cliapp.IntFlag{
				Name:  "port",
				Value: 26657,
				Usage: "this is the blockchain port",
			},
			cliapp.StringFlag{
				Name:  "dir",
				Value: "./blockchain",
				Usage: "this is the blockchain database path",
			},
			cliapp.StringFlag{
				Name:  "pass",
				Value: "ADck5qlB",
				Usage: "this is the password used to decrypt your configuration file",
			},
			cliapp.StringFlag{
				Name:  "file",
				Value: "xmn.conf",
				Usage: "this is the path of your encrypted configuration file",
			},
		},
		Action: func(c *cliapp.Context) error {
			// spawn the node:
			node := commands.SDKFunc.Spawn(commands.SpawnParams{
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

			// sleep 1 second before listening to keyboard:
			time.Sleep(time.Second * 1)
			termErr := term.Init()
			if termErr != nil {
				str := fmt.Sprintf("there was an error while enabling the keyboard listening: %s", termErr.Error())
				return errors.New(str)
			}
			defer term.Close()

			// blockchain started, loop until we stop:
			str := fmt.Sprintf("XMN cloud blockchain spawned, IP: %s\nPress Esc to stop...", client.IP())
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
	}
}
