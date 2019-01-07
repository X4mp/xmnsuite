package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveListFromWallet() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrievelist-fromwallet",
		Aliases: []string{"r"},
		Usage:   "Retrieves the requests that can be voted by users on the given walletid",
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
				Usage: "This is the wallet we want to retrieve requests from",
			},
			cliapp.IntFlag{
				Name:  "index",
				Value: 0,
				Usage: "The index of the list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of users to retrieve",
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

			activeRequestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the walletID:
			walletIDAsString := c.String("walletid")
			walletID, walletIDErr := uuid.FromString(walletIDAsString)
			if walletIDErr != nil {
				str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the from wallet:
			wal, walErr := walletRepository.RetrieveByID(&walletID)
			if walErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), walErr.Error())
				panic(errors.New(str))
			}

			// retrieve the requests:
			index := c.Int("index")
			amount := c.Int("amount")
			reqPS, reqPSErr := activeRequestRepository.RetrieveSetByWallet(wal, index, amount)
			if reqPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving active requests by wallet (ID: %s): %s", wal.ID().String(), reqPSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: reqPS,
			})

			// returns:
			return nil
		},
	}
}
