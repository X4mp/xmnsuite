package cli

import (
	"encoding/json"
	"fmt"
	"net"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/commands"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/helpers"
)

func retrieveGenesis() *cliapp.Command {
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
			helpers.Print(str)

			// returns:
			return nil
		},
	}
}
