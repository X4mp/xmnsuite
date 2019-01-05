package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"r"},
		Usage:   "Retrieves a transfer instance by id",
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
				Name:  "transferid",
				Value: "",
				Usage: "This is the transferid",
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

			transferRepository := transfer.SDKFunc.CreateRepository(transfer.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the transferid:
			transferIDAsString := c.String("transferid")
			trsfID, trsfIDErr := uuid.FromString(transferIDAsString)
			if trsfIDErr != nil {
				str := fmt.Sprintf("the given transferid (ID: %s) is not a valid id", transferIDAsString)
				panic(errors.New(str))
			}

			// retrieve the transfer:
			trsf, trsfErr := transferRepository.RetrieveByID(&trsfID)
			if trsfErr != nil {
				str := fmt.Sprintf("there was a problem while retrieving the transfer instance (ID: %s): %s", trsfID.String(), trsfErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: trsf,
			})

			// returns:
			return nil
		},
	}
}
