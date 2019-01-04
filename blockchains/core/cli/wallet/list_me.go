package wallet

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/configs"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func listMe() *cliapp.Command {
	return &cliapp.Command{
		Name:    "me",
		Aliases: []string{"c"},
		Usage:   "Show the wallets the current users have shares in.",
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

			// create the wallet repository:
			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
					PK: conf.WalletPK(),
					Client: tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
						IPAsString: c.String("host"),
					}),
					RoutePrefix: "",
				}),
			})

			// retrieve the users:
			pubKey := conf.WalletPK().PublicKey()
			usrPS, usrPSErr := userRepository.RetrieveSetByPubKey(pubKey, 0, -1)
			if usrPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet set (PubKey: %s): %s", pubKey.String(), usrPSErr.Error())
				panic(errors.New(str))
			}

			// render the list:
			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins:     usrPS.Instances(),
				Message: fmt.Sprintf("Index: %d - Amount: %d - TotalAmount: %d", usrPS.Index(), usrPS.Amount(), usrPS.TotalAmount()),
			})

			// returns:
			return nil
		},
	}
}
