package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	term "github.com/nsf/termbox-go"
	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	cliapp "github.com/urfave/cli"
	module_chain "github.com/xmnservices/xmnsuite/modules/chain"
)

func reset() {
	term.Sync()
}

func main() {
	app := cliapp.NewApp()
	app.Name = "xmnsuite"
	app.Usage = "Builds standalone blockchain applications using lua scripting"
	app.Flags = []cliapp.Flag{
		cliapp.StringFlag{
			Name:  "ccsize",
			Value: strconv.Itoa(120),
			Usage: "this is the lua call stack size",
		},
		cliapp.StringFlag{
			Name:  "rsize",
			Value: strconv.Itoa(120 * 20),
			Usage: "this is the lua registry size",
		},
		cliapp.StringFlag{
			Name:  "dbpath",
			Value: "./db",
			Usage: "this is the blockchain database path",
		},
		cliapp.StringFlag{
			Name:  "nodepk",
			Value: "",
			Usage: "this is the first blockchain node private key",
		},
		cliapp.StringFlag{
			Name:  "id",
			Value: uuid.NewV4().String(),
			Usage: "this is the blockchain instance id (UUID v.4)",
		},
		cliapp.StringFlag{
			Name:  "rpubkeys",
			Value: "",
			Usage: "these are the comma seperated root pub keys (that can write to every route on the blockchain)",
		},
		cliapp.StringFlag{
			Name:  "connector",
			Value: "",
			Usage: "this is the other blockchain that our blockchain will be able to connect to",
		},
		cliapp.StringFlag{
			Name:  "modules",
			Value: "",
			Usage: "these are the modules that are mandatory in order to run the given lua script",
		},
	}

	app.Commands = []cliapp.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate elements used in blockchain development",
			Subcommands: []cliapp.Command{
				{
					Name:  "pair",
					Usage: "generate a new PrivateKey/PublicKey pair",
					Action: func(c *cliapp.Context) error {
						pk := ed25519.GenPrivKey()
						str := fmt.Sprintf("Private Key: %s\nPublic Key:  %s", hex.EncodeToString(pk.Bytes()), hex.EncodeToString(pk.PubKey().Bytes()))
						print(str)
						return nil
					},
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "runs a blockchain application",
			Flags: []cliapp.Flag{
				cliapp.StringFlag{
					Name:  "ccsize",
					Value: strconv.Itoa(120),
					Usage: "this is the lua call stack size",
				},
				cliapp.StringFlag{
					Name:  "rsize",
					Value: strconv.Itoa(120 * 20),
					Usage: "this is the lua registry size",
				},
				cliapp.StringFlag{
					Name:  "dbpath",
					Value: "./db",
					Usage: "this is the blockchain database path",
				},
				cliapp.StringFlag{
					Name:  "nodepk",
					Value: "",
					Usage: "this is the first blockchain node private key",
				},
				cliapp.StringFlag{
					Name:  "id",
					Value: uuid.NewV4().String(),
					Usage: "this is the blockchain instance id (UUID v.4)",
				},
				cliapp.StringFlag{
					Name:  "rpubkeys",
					Value: "",
					Usage: "these are the comma seperated root pub keys (that can write to every route on the blockchain)",
				},
				cliapp.StringFlag{
					Name:  "connector",
					Value: "",
					Usage: "this is the other blockchain that our blockchain will be able to connect to",
				},
				cliapp.StringFlag{
					Name:  "modules",
					Value: "",
					Usage: "these are the modules that are mandatory in order to run the given lua script",
				},
			},
			Action: func(c *cliapp.Context) error {
				// retrieve the script path:
				scriptPath := c.Args().First()
				if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
					str := fmt.Sprintf("the given lua script path (%s) is invalid", scriptPath)
					return errors.New(str)
				}

				// create the cli instance:
				cli, cliErr := createCLI(c)
				if cliErr != nil {
					return cliErr
				}

				// execute the script:
				execErr := cli.execute(scriptPath)
				if execErr != nil {
					return execErr
				}

				// if there is a chain module, spawn:
				chainModule := cli.getModuleByName("chain")
				if chainModule != nil {
					node, nodeErr := chainModule.(module_chain.Chain).Spawn()
					if nodeErr != nil {
						// output error:
						str := fmt.Sprintf("there was an error while spawning a blockchain node: %s", nodeErr.Error())
						return errors.New(str)
					}
					defer node.Stop()
					node.Start()

					// sleep 5 seconds before listening to keyboard:
					time.Sleep(time.Second * 1)
					termErr := term.Init()
					if termErr != nil {
						str := fmt.Sprintf("there was an error while enabling the keyboard listening: %s", termErr.Error())
						return errors.New(str)
					}
					defer term.Close()

					// blockchain started, loop until we stop:
					print("Started... \nPress Esc to stop...")

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
				}

				return nil

			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func write(str string) string {
	out := fmt.Sprintf("\n************ xmnsuite ************\n")
	out = fmt.Sprintf("%s%s", out, str)
	out = fmt.Sprintf("%s\n********** end xmnsuite **********\n", out)
	return out
}

func print(str string) {
	fmt.Printf("%s", write(str))
}
