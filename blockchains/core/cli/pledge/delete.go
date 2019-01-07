package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func delete() *cliapp.Command {
	return &cliapp.Command{
		Name:    "delete",
		Aliases: []string{"s"},
		Usage:   "Delete creates a request to the wallet shareholders in order to delete a tokens pledge made to another wallet",
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
				Name:  "pledgeid",
				Value: "",
				Usage: "This is the id of the pledge to delete.",
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

			pledgeRepository := pledge.SDKFunc.CreateRepository(pledge.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the pledgeID:
			pledgeIDAsString := c.String("pledgeid")
			pledgeID, pledgeIDErr := uuid.FromString(pledgeIDAsString)
			if pledgeIDErr != nil {
				str := fmt.Sprintf("the given pledgeid (ID: %s) is not a valid id", pledgeIDAsString)
				panic(errors.New(str))
			}

			// retrieve the pledge:
			pldge, pldgeErr := pledgeRepository.RetrieveByID(&pledgeID)
			if pldgeErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the pledge (ID: %s): %s", pledgeID.String(), pldgeErr)
				panic(errors.New(str))
			}

			// save the request:
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: information.SDKFunc.CreateRepresentation(),
				DeleteEntity:         pldge,
			})

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: req,
			})

			// returns:
			return nil
		},
	}
}
