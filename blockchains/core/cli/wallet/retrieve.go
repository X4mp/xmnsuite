package wallet

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"s"},
		Usage:   "Retrieve retrieves a wallet by its id",
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
				Name:  "walletid",
				Value: "",
				Usage: "This is the id of the wallet to retrieve",
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

			// parse the walletID:
			walletIDAsString := c.String("walletid")
			walletID, walletIDErr := uuid.FromString(walletIDAsString)
			if walletIDErr != nil {
				str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the wallet:
			wal, walErr := walletReposiotry.RetrieveByID(&walletID)
			if walErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), walErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: wal,
			})

			// returns:
			return nil
		},
	}
}
