package active

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveList() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrievelist",
		Aliases: []string{"r"},
		Usage:   "Retrieves a list of requests",
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
				Usage: "The index of the list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of users to retrieve",
			},
			cliapp.StringFlag{
				Name:  "keyname",
				Value: "",
				Usage: "This keyname we want requests from",
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

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			activeRequestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the keyname:
			name := c.String("keyname")
			kname, knameErr := keynameRepository.RetrieveByName(name)
			if knameErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the keyname (Name: %s): %s", name, knameErr.Error())
				panic(errors.New(str))
			}

			// retrieve the requests:
			index := c.Int("index")
			amount := c.Int("amount")
			reqPS, reqPSErr := activeRequestRepository.RetrieveSetByKeyname(kname, index, amount)
			if reqPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving active requests: %s", reqPSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: reqPS,
			})

			// returns:
			return nil
		},
	}
}
