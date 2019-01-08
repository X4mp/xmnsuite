package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func create(met meta.Meta) *cliapp.Command {
	return &cliapp.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "Create creates a vote on a request",
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
				Usage: "This is the request id the vote will be binded on",
			},
			cliapp.StringFlag{
				Name:  "walletid",
				Value: "",
				Usage: "The walletid you want to vote with",
			},
			cliapp.StringFlag{
				Name:  "reason",
				Value: "",
				Usage: "The reason why you voted that way.  This will guide others in making their vote.",
			},
			cliapp.BoolFlag{
				Name:  "is_approved",
				Usage: "Whether the vote is approved or not",
			},
			cliapp.BoolFlag{
				Name:  "is_neutral",
				Usage: "True if you could not decide if the vote should be approved or not",
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

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			voteService := vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
				PK:          conf.WalletPK(),
				Client:      client,
				RoutePrefix: "",
			})

			activeVoteRepository := active_vote.SDKFunc.CreateRepository(active_vote.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the requestID:
			requestIDAsString := c.String("requestid")
			requestID, requestIDErr := uuid.FromString(requestIDAsString)
			if requestIDErr != nil {
				str := fmt.Sprintf("the given requestid (ID: %s) is not a valid id", requestIDAsString)
				panic(errors.New(str))
			}

			// retrieve request by its id:
			req, reqErr := activeRequestRepository.RetrieveByID(&requestID)
			if reqErr != nil {
				str := fmt.Sprintf("there was an error while retrieving a request (ID: %s): %s", requestID.String(), reqErr.Error())
				panic(errors.New(str))
			}

			// retrieve the user:
			fromUser, fromUserErr := userRepository.RetrieveByPubKey(conf.WalletPK().PublicKey())
			if fromUserErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the user (PubKey: %s): %s", conf.WalletPK().PublicKey().String(), fromUserErr.Error())
				panic(errors.New(str))
			}

			// create the vote request:
			vot := vote.SDKFunc.Create(vote.CreateParams{
				Request:    req,
				Voter:      fromUser,
				Reason:     c.String("reason"),
				IsApproved: c.Bool("is_approved"),
				IsNeutral:  c.Bool("is_neutral"),
			})

			// find the representation:
			kname := req.Request().Keyname()
			wr := met.WriteOnEntityRequest()
			reps := wr[kname.Group().Name()].Map()

			// save the vote:
			saveErr := voteService.Save(vot, reps[kname.Name()])
			if saveErr != nil {
				str := fmt.Sprintf("there was an error while saving a vote: %s", saveErr.Error())
				panic(errors.New(str))
			}

			// retrieve the active vote:
			activeVote, activeVoteErr := activeVoteRepository.RetrieveByVote(vot)
			if activeVoteErr != nil {
				str := fmt.Sprintf("there was an error while retrieving an active vote by vote (ID: %s): %s", vot.ID().String(), activeVoteErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: activeVote,
			})

			// returns:
			return nil
		},
	}
}
