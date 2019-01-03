package wallet

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/bitcoin/configs"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func create() *cliapp.Command {
	return &cliapp.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "Creates a new wallet and attach the current user has its first shareholder",
		Flags: []cliapp.Flag{
			cliapp.StringFlag{
				Name:  "id",
				Value: "",
				Usage: "This is the first user id of the new wallet (optional)",
			},
			cliapp.StringFlag{
				Name:  "shares",
				Value: "",
				Usage: "The amount of shares your current user will hold in the new wallet.",
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

			// retrieve the configurations:
			fileAsString := c.String("file")
			confRepository := configs.SDKFunc.CreateRepository()
			conf, confErr := confRepository.Retrieve(fileAsString, c.String("pass"))
			if confErr != nil {
				str := fmt.Sprintf("the given file (%s) either does not exist or the given password is invalid", fileAsString)
				panic(errors.New(str))
			}

			// process the request:
			pubKeyAsString := conf.WalletPK().PublicKey().String()
			req := helpers.SDKFunc.ProcessWalletRequest(helpers.ProcessWalletRequestParams{
				CLIContext:           c,
				EntityRepresentation: user.SDKFunc.CreateRepresentation(),
				Storable: user.SDKFunc.CreateNormalized(user.CreateNormalizedParams{
					ID:     c.String("id"),
					PubKey: pubKeyAsString,
					Shares: c.Int("shares"),
					Wallet: wallet.SDKFunc.CreateNormalized(wallet.CreateNormalizedParams{
						ID:              c.String("wallet_id"),
						CreatorPubKey:   pubKeyAsString,
						ConcensusNeeded: c.Int("concensus"),
					}),
				}),
			})

			helpers.SDKFunc.PrintSuccessNewInstance(helpers.PrintSuccessNewInstanceParams{
				Ins:     req,
				Message: "Success!  The wallet request has been saved.",
			})

			// returns:
			return nil
		},
	}
}
