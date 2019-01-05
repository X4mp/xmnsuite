package information

import (
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func update() *cliapp.Command {
	return &cliapp.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Creates a request to the token holders in order to update the information instance",
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
				Name:  "gazprice",
				Value: 0,
				Usage: "This is the gaz price paid per kb when executing transactions",
			},
			cliapp.IntFlag{
				Name:  "concensus",
				Value: 0,
				Usage: "This is the amount of token that needs to vote in a direction in order to have concensus",
			},
			cliapp.IntFlag{
				Name:  "max_validators",
				Value: 0,
				Usage: "This is the maximum amount of validators the blockchain can have",
			},
			cliapp.StringFlag{
				Name:  "reason",
				Value: "",
				Usage: "This is the reason why you want to update the information instance",
			},
			cliapp.StringFlag{
				Name:  "walletid",
				Value: "",
				Usage: "This is the walletID from which you want to send the request",
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

			// create the information repository:
			infRepository := information.SDKFunc.CreateRepository(information.CreateRepositoryParams{
				EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
					PK:     conf.WalletPK(),
					Client: client,
				}),
			})

			// retrieve the information:
			inf, infErr := infRepository.Retrieve()
			if infErr != nil {
				panic(infErr)
			}

			newGazPrice := c.Int("gazprice")
			if newGazPrice == 0 {
				newGazPrice = inf.GazPricePerKb()
			}

			newConcensusNeeded := c.Int("concensus")
			if newConcensusNeeded == 0 {
				newConcensusNeeded = inf.ConcensusNeeded()
			}

			maxAmountValidators := c.Int("max_validators")
			if maxAmountValidators == 0 {
				maxAmountValidators = inf.MaxAmountOfValidators()
			}

			// save the request:
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: information.SDKFunc.CreateRepresentation(),
				Ins: information.SDKFunc.Create(information.CreateParams{
					ID:                    inf.ID(),
					GazPricePerKb:         newGazPrice,
					ConcensusNeeded:       newConcensusNeeded,
					MaxAmountOfValidators: maxAmountValidators,
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
