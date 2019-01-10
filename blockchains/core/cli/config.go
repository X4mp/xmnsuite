package cli

import (
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/commands"
)

func generateConfig() *cliapp.Command {
	return &cliapp.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Generate the config file",
		Flags: []cliapp.Flag{
			cliapp.StringFlag{
				Name:  "file",
				Value: "xmn.conf",
				Usage: "this is the path of your encrypted configuration file",
			},
			cliapp.StringFlag{
				Name:  "pass",
				Value: "ADck5qlB",
				Usage: "this is the password used to decrypt your configuration file",
			},
			cliapp.StringFlag{
				Name:  "retypePass",
				Value: "ADck5qlB",
				Usage: "this is the retyped password used to decrypt your configuration file",
			},

		},
		Action: func(c *cliapp.Context) error {
			_ = commands.SDKFunc.GenerateConfigs(commands.GenerateConfigsParams{
				Filename: c.String("file"),
				Pass: c.String("pass"),
				RetypedPass: c.String("retypePass"),
			})

			return nil
		},
	}
}

