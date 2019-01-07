package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"r"},
		Usage:   "Retrieve retrieves a vote by id",
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
				Name:  "voteid",
				Value: "",
				Usage: "The id of the vote to delete",
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

			activeVoteRepository := active_vote.SDKFunc.CreateRepository(active_vote.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the voteid:
			voteIDAsString := c.String("voteid")
			voteID, voteIDErr := uuid.FromString(voteIDAsString)
			if voteIDErr != nil {
				str := fmt.Sprintf("the given voteid (ID: %s) is not a valid id", voteIDAsString)
				panic(errors.New(str))
			}

			// retrieve the vote:
			vot, votErr := activeVoteRepository.RetrieveByID(&voteID)
			if votErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the vote (ID: %s): %s", voteID.String(), votErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: vot,
			})

			// returns:
			return nil
		},
	}
}
