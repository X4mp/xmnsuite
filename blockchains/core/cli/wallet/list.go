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

func list() *cliapp.Command {
	return &cliapp.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "Creates a new wallet and attach the current user has its first shareholder",
		Flags: []cliapp.Flag{
			cliapp.IntFlag{
				Name:  "index",
				Value: 0,
				Usage: "This is the list index",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "This is the amount of elements to retrieve",
			},
			cliapp.StringFlag{
				Name:  "wallet_id",
				Value: "",
				Usage: "This is the wallet id (optional)",
			},
			cliapp.StringFlag{
				Name:  "concensus",
				Value: "",
				Usage: "This is the amount of shares that needs to vote in order to approve or disapprove a wallet request",
			},
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

			// create the wallet repository:
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{}),
			})

			// retrieve wallets:
			walPS, walPSErr := walletRepository.RetrieveSet(c.Int("index"), c.Int("amount"))
			if walPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving a wallet list: %s", walPSErr.Error())
				panic(errors.New(str))
			}

			// render the list:
			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins:     walPS.Instances(),
				Message: fmt.Sprintf("Index: %d - Amount: %d - TotalAmount: %d", walPS.Index(), walPS.Amount(), walPS.TotalAmount()),
			})

			// returns:
			return nil
		},
	}
}
