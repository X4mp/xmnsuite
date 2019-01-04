package genesis

import (
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"r"},
		Usage:   "Retrieve the genesis instance",
		Flags: []cliapp.Flag{
			cliapp.StringFlag{
				Name:  "host",
				Value: "",
				Usage: "This is the blockchain ip and port (example: 127.0.0.1:8080)",
			},
			cliapp.StringFlag{
				Name:  "file",
				Value: "",
				Usage: "This is the path of your encrypted configuration file",
			},
			cliapp.StringFlag{
				Name:  "pass",
				Value: "",
				Usage: "This is the password used to decrypt the encrypted configuration file",
			},
		},
		Action: func(c *cliapp.Context) error {
			defer func() {
				if r := recover(); r != nil {
					str := fmt.Sprintf("%s", r)
					core_helpers.Print(str)
				}
			}()

			// retrieve conf with client:
			conf, client := helpers.SDKFunc.RetrieveConfWithClient(helpers.RetrieveConfWithClientParams{
				CLIContext: c,
			})

			// create the genesis repository:
			genRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
					PK:     conf.WalletPK(),
					Client: client,
				}),
			})

			// retrieve the genesis:
			gen, genErr := genRepository.Retrieve()
			if genErr != nil {
				panic(genErr)
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins:     gen,
				Message: "Success!  Here's the Genesis instance:",
			})

			// returns:
			return nil
		},
	}
}
