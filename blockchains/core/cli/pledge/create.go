package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func create() *cliapp.Command {
	return &cliapp.Command{
		Name:    "create",
		Aliases: []string{"s"},
		Usage:   "Create creates a request to the wallet shareholders in order to create a tokens pledge to another wallet",
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
				Usage: "This is the from walletid.  The request will be voted by the shareholders of that wallet.",
			},
			cliapp.StringFlag{
				Name:  "towalletid",
				Value: "",
				Usage: "This is the to walletid",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 0,
				Usage: "This is the amount of tokens to pledge",
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

			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the genesis:
			gen, genErr := genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
				panic(errors.New(str))
			}

			// parse the fromWalletID:
			fromWalletIDAsString := c.String("fromwalletid")
			fromWalletID, fromWalletIDErr := uuid.FromString(fromWalletIDAsString)
			if fromWalletIDErr != nil {
				str := fmt.Sprintf("the given fromwalletid (ID: %s) is not a valid id", fromWalletIDAsString)
				panic(errors.New(str))
			}

			// parse the toWalletID:
			toWalletIDAsString := c.String("towalletid")
			toWalletID, toWalletIDErr := uuid.FromString(toWalletIDAsString)
			if toWalletIDErr != nil {
				str := fmt.Sprintf("the given towalletid (ID: %s) is not a valid id", toWalletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the from wallet:
			fromWallet, fromWalletErr := walletRepository.RetrieveByID(&fromWalletID)
			if fromWalletErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", fromWalletID.String(), fromWalletErr)
				panic(errors.New(str))
			}

			// retrieve the to wallet:
			toWallet, toWalletErr := walletRepository.RetrieveByID(&toWalletID)
			if toWalletErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", toWalletID.String(), toWalletErr)
				panic(errors.New(str))
			}

			// save the request:
			amount := c.Int("amount")
			tok := gen.Deposit().Token()
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: information.SDKFunc.CreateRepresentation(),
				SaveEntity: pledge.SDKFunc.Create(pledge.CreateParams{
					From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
						From:   fromWallet,
						Token:  tok,
						Amount: amount,
					}),
					To: toWallet,
				}),
			})

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: req,
			})

			// returns:
			return nil
		},
	}
}
