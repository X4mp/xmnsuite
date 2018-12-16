package cli

import (
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/applications/forex/commands"
	"github.com/xmnservices/xmnsuite/helpers"
)

func generateConfigs() *cliapp.Command {
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
			helpers.Print(str)

			// returns:
			return nil
		},
	}
}
