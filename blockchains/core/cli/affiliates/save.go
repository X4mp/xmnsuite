package affiliates

import (
	"errors"
	"fmt"
	"net/url"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func save() *cliapp.Command {
	return &cliapp.Command{
		Name:    "save",
		Aliases: []string{"s"},
		Usage:   "Saves an affiliate instance",
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
				Name:  "walletid",
				Value: "",
				Usage: "This is the walletid associated with the affiliate",
			},
			cliapp.StringFlag{
				Name:  "url",
				Value: "",
				Usage: "This is the url of the affiliate web service",
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

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// parse the url:
			rawURL := c.String("url")
			url, urlErr := url.Parse(rawURL)
			if urlErr != nil {
				str := fmt.Sprintf("there was an error while parsing the raw URL (%s): %s", rawURL, urlErr.Error())
				panic(errors.New(str))
			}

			// parse the walletid:
			walletIDAsString := c.String("walletid")
			walID, walIDErr := uuid.FromString(walletIDAsString)
			if walIDErr != nil {
				str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
				panic(errors.New(str))
			}

			// retrieve the wallet:
			wal, walErr := walletRepository.RetrieveByID(&walID)
			if walErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walID.String(), walErr.Error())
				panic(errors.New(str))
			}

			// create the request:
			req := helpers.SDKFunc.SaveRequest(helpers.SaveRequestParams{
				CLIContext:           c,
				EntityRepresentation: affiliates.SDKFunc.CreateRepresentation(),
				SaveEntity: affiliates.SDKFunc.Create(affiliates.CreateParams{
					Owner: wal,
					URL:   url,
				}),
			})

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: req,
			})

			// returns:
			return nil
		},
	}
}
