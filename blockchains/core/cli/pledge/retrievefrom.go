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

func retrieveFrom() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-from",
		Aliases: []string{"s"},
		Usage:   "Retrieves a list of pledge made from a given walletid",
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
				Name:  "fromwalletid",
				Value: "",
				Usage: "The id that we want to retrieve pledges made from this wallet",
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

			// parse the fromWalletID:
			fromWalletIDAsString := c.String("fromwalletid")
			fromWalletID, fromWalletIDErr := uuid.FromString(fromWalletIDAsString)
			if fromWalletIDErr != nil {
				str := fmt.Sprintf("the given fromwalletid (ID: %s) is not a valid id", fromWalletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the from wallet:
			fromWallet, fromWalletErr := walletRepository.RetrieveByID(&fromWalletID)
			if fromWalletErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the fromWallet (ID: %s): %s", fromWalletID.String(), fromWalletErr.Error())
				panic(errors.New(str))
			}

			// retrieve the pledge set:
			index := c.Int("index")
			amount := c.Int("amount")
			pledgePS, pledgePSErr := pledgeRepository.RetrieveSetByFromWallet(fromWallet, index, amount)
			if pledgePSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the pledge set by fromWallet (ID: %s): %s", fromWallet.ID().String(), pledgePSErr.Error())
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
