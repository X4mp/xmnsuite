package affiliates

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"r"},
		Usage:   "Retrieve retrieves an affiliate by its id",
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
				Name:  "affid",
				Value: "",
				Usage: "This is the id of the affiliate",
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

			affiliateRepository := affiliates.SDKFunc.CreateRepository(affiliates.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the affID:
			affIDAsString := c.String("affid")
			affID, affIDErr := uuid.FromString(affIDAsString)
			if affIDErr != nil {
				str := fmt.Sprintf("the given affid (ID: %s) is not a valid id", affIDAsString)
				panic(errors.New(str))
			}

			// retrieve the affiliate:
			aff, affErr := affiliateRepository.RetrieveByID(&affID)
			if affErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the affiliate (ID: %s): %s", affID.String(), affErr.Error())
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
