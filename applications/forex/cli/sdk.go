package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/applications/forex/commands"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
)

func reset() {
	term.Sync()
}

// SDKFunc represents the CLI sdk func
var SDKFunc = struct {
	GenerateConfigs func() *cliapp.Command
	SpawnMain       func() *cliapp.Command
	RetrieveGenesis func() *cliapp.Command
}{
	GenerateConfigs: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generates the encrypted configuration file (private keys) in order to interact with the XMN forex blockchain",
			Flags: []cliapp.Flag{
				cliapp.StringFlag{
					Name:  "pass",
					Value: "",
					Usage: "this is the password used to decrypt your configuration file",
				},
				cliapp.StringFlag{
					Name:  "rpass",
					Value: "",
					Usage: "this is the retyped pass... same as pass, to make sure you typed it correctly",
				},
				cliapp.StringFlag{
					Name:  "file",
					Value: "",
					Usage: "this is the path of your encrypted configuration file",
				},
			},
			Action: func(c *cliapp.Context) error {
				// generate the config file:
				commands.SDKFunc.GenerateConfigs(commands.GenerateConfigsParams{
					Pass:        c.String("pass"),
					RetypedPass: c.String("rpass"),
					Filename:    c.String("file"),
				})

				str := fmt.Sprintf("Successful!  Ecncrypted configuration file saved: %s", c.String("file"))
				print(str)

				// returns:
				return nil
			},
		}
	},
	SpawnMain: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "spawn",
			Aliases: []string{"s"},
			Usage:   "Spawns the main forex blockchain",
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
				// spawn the node:
				node := commands.SDKFunc.SpawnMain(commands.SpawnMainParams{
					Pass:     c.String("pass"),
					Filename: c.String("file"),
					Dir:      c.String("dir"),
					Port:     c.Int("port"),
				})

				// render to the cli:
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
				str := fmt.Sprintf("XMN main blockchain spawned, IP: %s\nPress Esc to stop...", client.IP())
				print(str)

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
	},
	RetrieveGenesis: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "genesis",
			Aliases: []string{"rg"},
			Usage:   "Retrieves the genesis transaction of the blockchain",
			Flags: []cliapp.Flag{
				cliapp.IntFlag{
					Name:  "port",
					Value: 26657,
					Usage: "this is the blockchain port",
				},
				cliapp.StringFlag{
					Name:  "ip",
					Value: "127.0.0.1",
					Usage: "this is the blockchain ip address",
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

				fmt.Printf("\n\n Pass: %s \n\n", c.String("pass"))
				fmt.Printf("\n\n File: %s \n\n", c.String("file"))
				fmt.Printf("\n\n IP: %s \n\n", c.String("ip"))
				fmt.Printf("\n\n Port: %d \n\n", c.Int("port"))

				// retrieve:
				gen := commands.SDKFunc.RetrieveGenesis(commands.RetrieveGenesisParams{
					Pass:     c.String("pass"),
					Filename: c.String("file"),
					IP:       net.ParseIP(c.String("ip")),
					Port:     c.Int("port"),
				})

				// normalize:
				normalized, normalizedErr := genesis.SDKFunc.CreateMetaData().Normalize()(gen)
				if normalizedErr != nil {
					panic(normalizedErr)
				}

				// beauty-print:
				data, dataErr := json.MarshalIndent(normalized, "", "    ")
				if dataErr != nil {
					panic(dataErr)
				}

				str := fmt.Sprintf(string(data))
				print(str)

				// returns:
				return nil
			},
		}
	},
}

func write(str string) string {
	out := fmt.Sprintf("\n************ XMN Forex Exchange ************\n")
	out = fmt.Sprintf("%s%s", out, str)
	out = fmt.Sprintf("%s\n********** end XMN Forex Exchange **********\n", out)
	return out
}

func print(str string) {
	fmt.Printf("%s", write(str))
}