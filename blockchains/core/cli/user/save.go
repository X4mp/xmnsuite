package user

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
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
		Usage:   "Save creates a new user on a new wallet",
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
				Name:  "shares",
				Value: 100,
				Usage: "This is the amount of shares",
			},
			cliapp.StringFlag{
				Name:  "name",
				Value: "",
				Usage: "The name of the user.  Accepted characters are letters, numbers and underscores (_).  Minimum length: 3 characters.",
			},
			cliapp.IntFlag{
				Name:  "cneeded",
				Value: 100,
				Usage: "This is the amount of shares that needs to vote in a direction in order to reach concensus on a decision",
			},
			cliapp.StringFlag{
				Name:  "pubkey",
				Value: "",
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
			conf, _ := helpers.SDKFunc.RetrieveConfWithClient(helpers.RetrieveConfWithClientParams{
				CLIContext: c,
			})

			// create the pubKey:
			pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
				PubKeyAsString: pubKeyAsString,
			})

			// save the request:
			shares := c.Int("shares")
			concensusNeeded := c.Int("cneeded")
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: information.SDKFunc.CreateRepresentation(),
				SaveEntity: user.SDKFunc.Create(user.CreateParams{
					Name:   c.String("name"),
					PubKey: pubKey,
					Shares: shares,
					Wallet: wallet.SDKFunc.Create(wallet.CreateParams{
						Creator:         conf.WalletPK().PublicKey(),
						ConcensusNeeded: concensusNeeded,
					}),
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
