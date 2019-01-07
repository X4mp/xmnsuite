package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieveList() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-list",
		Aliases: []string{"r"},
		Usage:   "Retrieve retrieves a list of votes related to a request",
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
				Name:  "requestid",
				Value: "",
				Usage: "The requestid of the votes",
			},
			cliapp.BoolFlag{
				Name:  "is_approved",
				Usage: "Whether the votes are approved or not",
			},
			cliapp.BoolFlag{
				Name:  "is_neutral",
				Usage: "Whether the votes are neutral",
			},
			cliapp.IntFlag{
				Name:  "index",
				Value: 0,
				Usage: "The index of the list",
			},
			cliapp.IntFlag{
				Name:  "amount",
				Value: 20,
				Usage: "The amount of votes to retrieve",
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

			activeRequestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			activeVoteRepository := active_vote.SDKFunc.CreateRepository(active_vote.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the requestid:
			requestIDAsString := c.String("requestid")
			requestID, requestIDErr := uuid.FromString(requestIDAsString)
			if requestIDErr != nil {
				str := fmt.Sprintf("the given requestid (ID: %s) is not a valid id", requestIDAsString)
				panic(errors.New(str))
			}

			// retrieve the request:
			req, reqErr := activeRequestRepository.RetrieveByID(&requestID)
			if reqErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the request (ID: %s): %s", requestID.String(), reqErr.Error())
				panic(errors.New(str))
			}

			// retrieve the votes:
			index := c.Int("index")
			amount := c.Int("amount")
			isApproved := c.Bool("is_approved")
			isNeutral := c.Bool("is_neutral")
			votesPS, votesPSErr := activeVoteRepository.RetrieveSetByRequestWithDirection(req, index, amount, isApproved, isNeutral)
			if votesPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving votes set by request (ID: %s): %s", req.ID().String(), votesPSErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: votesPS,
			})

			// returns:
			return nil
		},
	}
}
