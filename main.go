package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	term "github.com/nsf/termbox-go"
	uuid "github.com/satori/go.uuid"
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	cli "github.com/urfave/cli"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	xmnmodule "github.com/xmnservices/xmnsuite/modules/xmn"
	"github.com/xmnservices/xmnsuite/tendermint"
	lua "github.com/yuin/gopher-lua"
)

var cdc = amino.NewCodec()

func init() {
	// crypto.PrivKey
	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.Ed25519PrivKeyAminoRoute, nil)
}

func reset() {
	term.Sync()
}

func main() {

	app := cli.NewApp()
	app.Name = "xmnsuite"
	app.Usage = "Builds standalone blockchain applications using lua scripting"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ccsize",
			Value: strconv.Itoa(120),
			Usage: "this is the lua call stack size",
		},
		cli.StringFlag{
			Name:  "rsize",
			Value: strconv.Itoa(120 * 20),
			Usage: "this is the lua registry size",
		},
		cli.StringFlag{
			Name:  "dbpath",
			Value: "./db",
			Usage: "this is the blockchain database path",
		},
		cli.StringFlag{
			Name:  "nodepk",
			Value: "",
			Usage: "this is the first blockchain node private key",
		},
		cli.StringFlag{
			Name:  "id",
			Value: uuid.NewV4().String(),
			Usage: "this is the blockchain instance id (UUID v.4)",
		},
		cli.StringFlag{
			Name:  "rpubkeys",
			Value: "",
			Usage: "these are the comma seperated root pub keys (that can write to every route on the blockchain)",
		},
		cli.StringFlag{
			Name:  "connector",
			Value: "",
			Usage: "this is the other blockchain that our blockchain will be able to connect to",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate elements used in blockchain development",
			Subcommands: []cli.Command{
				{
					Name:  "pair",
					Usage: "generate a new PrivateKey/PublicKey pair",
					Action: func(c *cli.Context) error {
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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ccsize",
					Value: strconv.Itoa(120),
					Usage: "this is the lua call stack size",
				},
				cli.StringFlag{
					Name:  "rsize",
					Value: strconv.Itoa(120 * 20),
					Usage: "this is the lua registry size",
				},
				cli.StringFlag{
					Name:  "dbpath",
					Value: "./db",
					Usage: "this is the blockchain database path",
				},
				cli.StringFlag{
					Name:  "nodepk",
					Value: "",
					Usage: "this is the first blockchain node private key",
				},
				cli.StringFlag{
					Name:  "id",
					Value: uuid.NewV4().String(),
					Usage: "this is the blockchain instance id (UUID v.4)",
				},
				cli.StringFlag{
					Name:  "rpubkeys",
					Value: "",
					Usage: "these are the comma seperated root pub keys (that can write to every route on the blockchain)",
				},
				cli.StringFlag{
					Name:  "connector",
					Value: "",
					Usage: "this is the other blockchain that our blockchain will be able to connect to",
				},
			},
			Action: func(c *cli.Context) error {

				termErr := term.Init()
				if termErr != nil {
					str := fmt.Sprintf("there was an error while enabling the keyboard listening: %s", termErr.Error())
					return errors.New(str)
				}
				defer term.Close()

				scriptPath := c.Args().First()
				if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
					str := fmt.Sprintf("the given lua script path (%s) is invalid", scriptPath)
					return errors.New(str)
				}

				// ccsize:
				ccSizeAsString := c.String("ccsize")
				ccSize, ccSizeErr := strconv.Atoi(ccSizeAsString)
				if ccSizeErr != nil {
					// log:
					log.Printf("there was an error while converting a string to an int: %s", ccSizeErr.Error())

					// output error:
					str := fmt.Sprintf("the ccsize param (%s) must be an int", ccSizeAsString)
					return errors.New(str)
				}

				// rsize:
				rSizeAsString := c.String("rsize")
				rSize, rSizeErr := strconv.Atoi(rSizeAsString)
				if rSizeErr != nil {
					// log:
					log.Printf("there was an error while converting a string to an int: %s", rSizeErr.Error())

					// output error:
					str := fmt.Sprintf("the rsize param (%s) must be an int", rSizeAsString)
					return errors.New(str)
				}

				// dbpath:
				dbPath := c.String("dbpath")

				// nodepk:
				nodePkAsString := c.String("nodepk")
				nodePKAsBytes, nodePKAsBytesErr := hex.DecodeString(nodePkAsString)
				if nodePKAsBytesErr != nil {
					// log:
					log.Printf("there was an error while decoding a string to hex: %s", nodePKAsBytesErr.Error())

					// output error:
					str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
					return errors.New(str)
				}

				nodePK := new(ed25519.PrivKeyEd25519)
				nodePKErr := cdc.UnmarshalBinaryBare(nodePKAsBytes, nodePK)
				if nodePKErr != nil {
					// log:
					log.Printf("there was an error while Unmarshalling []byte to PrivateKey instance: %s", nodePKErr.Error())

					// output error:
					str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
					return errors.New(str)
				}

				// id:
				idAsString := c.String("id")
				id, idErr := uuid.FromString(idAsString)
				if idErr != nil {
					// log:
					log.Printf("there was an error while converting a string to an ID: %s", idErr.Error())

					// output error:
					str := fmt.Sprintf("the given id (%s) is not a valid ID", idAsString)
					return errors.New(str)
				}

				// rpubkeys:

				// create the lua context:
				context := lua.NewState(lua.Options{
					CallStackSize: ccSize,
					RegistrySize:  rSize,
				})
				defer context.Close()

				// create the datastore:
				ds := datastore.SDKFunc.Create()

				// create the params:
				xmnParams := xmnmodule.ExecuteParams{
					DBPath:     dbPath,
					NodePK:     nodePK,
					InstanceID: &id,
					Store:      ds,
					Context:    context,
					ScriptPath: scriptPath,
					Client:     nil,
				}

				// connect to:
				connectorAsString := c.String("connector")
				if connectorAsString != "" {
					appService := tendermint.SDKFunc.CreateApplicationService()
					client, clientErr := appService.Connect(connectorAsString)
					if clientErr != nil {
						// log:
						log.Printf("there was an error while connecting to the given host: %s", clientErr.Error())

						// output error:
						str := fmt.Sprintf("the given connector (%s) is not a valid blockchain host", connectorAsString)
						return errors.New(str)
					}

					xmnParams.Client = client
				}

				// create XMN:
				xmnNode := xmnmodule.SDKFunc.Execute(xmnParams)
				defer xmnNode.Stop()

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
