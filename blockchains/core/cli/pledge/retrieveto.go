package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveTo() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-to",
		Aliases: []string{"s"},
		Usage:   "Retrieves a list of pledge made to a given walletid",
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
			cliapp.StringFlag{
				Name:  "towalletid",
				Value: "",
				Usage: "The id that we want to retrieve pledges made to this wallet",
			},
			cliapp.IntFlag{
				Name:  "index",
				Value: 0,
				Usage: "The index of the pledge list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of pledge to retrieve",
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

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			pledgeRepository := pledge.SDKFunc.CreateRepository(pledge.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the toWalletID:
			toWalletIDAsString := c.String("towalletid")
			toWalletID, toWalletIDErr := uuid.FromString(toWalletIDAsString)
			if toWalletIDErr != nil {
				str := fmt.Sprintf("the given towalletid (ID: %s) is not a valid id", toWalletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the to wallet:
			toWallet, toWalletErr := walletRepository.RetrieveByID(&toWalletID)
			if toWalletErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the toWallet (ID: %s): %s", toWalletID.String(), toWalletErr.Error())
				panic(errors.New(str))
			}

			// retrieve the pledge set:
			index := c.Int("index")
			amount := c.Int("amount")
			pledgePS, pledgePSErr := pledgeRepository.RetrieveSetByToWallet(toWallet, index, amount)
			if pledgePSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the pledge set by toWallet (ID: %s): %s", toWallet.ID().String(), pledgePSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: pledgePS,
			})

			// returns:
			return nil
		},
	}
}
