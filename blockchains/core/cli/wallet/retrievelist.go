package wallet

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveList() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-list",
		Aliases: []string{"s"},
		Usage:   "Retrieve retrieves a list of wallets",
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
			cliapp.IntFlag{
				Name:  "index",
				Value: 0,
				Usage: "The index of the wallet list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of wallet to retrieve",
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

			// create the repositories:
			entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
				PK:     conf.WalletPK(),
				Client: client,
			})

			walletReposiotry := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the validator set:
			index := c.Int("index")
			amount := c.Int("amount")
			walPS, walPSErr := walletReposiotry.RetrieveSet(index, amount)
			if walPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet set (index: %d, amount: %d): %s", index, amount, walPSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: walPS,
			})

			// returns:
			return nil
		},
	}
}
