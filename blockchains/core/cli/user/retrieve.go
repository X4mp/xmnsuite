package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"rl"},
		Usage:   "Retrieves a user by id",
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
				Name:  "userid",
				Value: "",
				Usage: "This is the id of the user to retrieve",
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

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the userID:
			userIDAsString := c.String("userid")
			userID, userIDErr := uuid.FromString(userIDAsString)
			if userIDErr != nil {
				str := fmt.Sprintf("the given userid (ID: %s) is not a valid id", userIDAsString)
				panic(errors.New(str))
			}

			// retrieve the user:
			usr, usrErr := userRepository.RetrieveByID(&userID)
			if usrErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the user (ID: %s): %s", userID.String(), usrErr)
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: usr,
			})

			// returns:
			return nil
		},
	}
}
