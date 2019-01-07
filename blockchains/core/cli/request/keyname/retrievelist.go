package keyname

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveList() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrievelist",
		Aliases: []string{"r"},
		Usage:   "Retrieves a list of keynames by group name",
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
				Name:  "groupname",
				Value: "",
				Usage: "The name of the group",
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

			groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the group:
			name := c.String("groupname")
			grp, grpErr := groupRepository.RetrieveByName(name)
			if grpErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the group (Name: %s): %s", name, grpErr.Error())
				panic(errors.New(str))
			}

			// retrieve the keynames:
			knamePS, knamePSErr := keynameRepository.RetrieveSetByGroup(grp, 0, -1)
			if knamePSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the keyname set by group (Name: %s): %s", name, knamePSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: knamePS,
			})

			// returns:
			return nil
		},
	}
}
