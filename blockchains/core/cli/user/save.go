package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/crypto"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func save() *cliapp.Command {
	return &cliapp.Command{
		Name:    "save",
		Aliases: []string{"c"},
		Usage:   "Save creates a request to the wallet shareholders in order to create/update a user on a given wallet",
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
				Usage: "This is the id of the wallet we want to add a new user on.  The request will also be sent from your user that is attached to this wallet.  If you have no user on that wallet, it returns an error.",
			},
			cliapp.IntFlag{
				Name:  "shares",
				Value: 0,
				Usage: "This is the amount of shares",
			},
			cliapp.IntFlag{
				Name:  "pubkey",
				Value: 0,
				Usage: "This is the new user public key",
			},
		},
		Action: func(c *cliapp.Context) error {
			defer func() {
				if r := recover(); r != nil {
					str := fmt.Sprintf("%s", r)
					core_helpers.Print(str)
				}
			}()

			pubKeyAsString := c.String("pubkey")
			if pubKeyAsString == "" {
				panic(errors.New("the pubKey cannot be empty"))
			}

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

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
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
			wal, walErr := walletRepository.RetrieveByID(&walletID)
			if walErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), walErr)
				panic(errors.New(str))
			}

			// create the pubKey:
			pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
				PubKeyAsString: pubKeyAsString,
			})

			// retrieve the user if the shares are not set:
			shares := c.Int("shares")
			if shares == 0 {
				retUser, retUserErr := userRepository.RetrieveByPubKeyAndWallet(pubKey, wal)
				if retUserErr != nil {
					str := fmt.Sprintf("the user (pubKey: %s, walletID: %s) does not exists, therefore the shares parameter is mandatory", pubKey.String(), wal.ID().String())
					panic(errors.New(str))
				}

				shares = retUser.Shares()
			}

			// save the request:
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: information.SDKFunc.CreateRepresentation(),
				SaveEntity: user.SDKFunc.Create(user.CreateParams{
					PubKey: pubKey,
					Shares: shares,
					Wallet: wal,
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
