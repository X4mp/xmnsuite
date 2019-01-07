package affiliates

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveByWallet() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-bywallet",
		Aliases: []string{"rw"},
		Usage:   "Retrieve retrieves an affiliate by a walletid",
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

			affiliateRepository := affiliates.SDKFunc.CreateRepository(affiliates.CreateRepositoryParams{
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

			// retrieve the affiliate:
			aff, affErr := affiliateRepository.RetrieveByWallet(wal)
			if affErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the affiliate (walletID: %s): %s", wal.ID().String(), affErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: aff,
			})

			// returns:
			return nil
		},
	}
}
