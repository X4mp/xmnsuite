package validator

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveList() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-list",
		Aliases: []string{"s"},
		Usage:   "Retrieve retrieves a list of validators",
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
				Name:  "index",
				Value: 0,
				Usage: "The index of the validator list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of validator to retrieve",
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

			validatorRepository := validator.SDKFunc.CreateRepository(validator.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the validator set:
			index := c.Int("index")
			amount := c.Int("amount")
			valPS, valPSErr := validatorRepository.RetrieveSet(index, amount)
			if valPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the validator set (index: %d, amount: %d): %s", index, amount, valPSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: valPS,
			})

			// returns:
			return nil
		},
	}
}
