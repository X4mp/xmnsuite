package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveByPledge() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-bypledge",
		Aliases: []string{"s"},
		Usage:   "Retrieve retrieves a validator instance by pledgeid",
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
				Usage: "This is the pledgeid associated with the validator",
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

			validatorRepository := validator.SDKFunc.CreateRepository(validator.CreateRepositoryParams{
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

			// retrieve the validator by pledge:
			val, valErr := validatorRepository.RetrieveByPledge(pldge)
			if valErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the validator instance by pledge (ID: %s): %s", pldge.ID().String(), valErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: val,
			})

			// returns:
			return nil
		},
	}
}
