package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/configs"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrievelist() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-list",
		Aliases: []string{"rl"},
		Usage:   "Retrieves a list of users related to a wallet",
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
				Usage: "This is the walletid we want to retrieve users from. If not set, all users are retrieved. (Optional)",
			},
			cliapp.BoolFlag{
				Name:  "me",
				Usage: "If set to true, retrieve only my users.",
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

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// set the index and amount:
			retUserPS, retUser, retUserPSErr := retrieveRightUsers(c, conf, userRepository, walletRepository)
			if retUserPSErr != nil {
				panic(retUserPSErr)
			}

			if retUserPS != nil {
				helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
					Ins: retUserPS,
				})

				return nil
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: retUser,
			})

			// returns:
			return nil
		},
	}
}

func retrieveRightUsers(c *cliapp.Context, conf configs.Configs, userRepository user.Repository, walletRepository wallet.Repository) (entity.PartialSet, entity.Entity, error) {
	// get the variables:
	me := c.Bool("me")
	index := c.Int("index")
	amount := c.Int("amount")
	walletIDAsString := c.String("walletid")

	// if there is no wallet:
	if walletIDAsString == "" {
		if !me {
			retUsersPS, retUsersPSErr := userRepository.RetrieveSet(index, amount)
			if retUsersPSErr != nil {
				str := fmt.Sprintf("there was an error while retrieving a set of users (index: %d, amount: %d): %s", index, amount, retUsersPSErr.Error())
				return nil, nil, errors.New(str)
			}

			return retUsersPS, nil, nil
		}

		retUser, retUserErr := userRepository.RetrieveByPubKey(conf.WalletPK().PublicKey())
		if retUserErr != nil {
			str := fmt.Sprintf("there was an error while retrieving a user (pubKey: %s): %s", conf.WalletPK().PublicKey().String(), retUserErr.Error())
			panic(errors.New(str))
		}

		return nil, retUser, nil
	}

	// parse the walletID:
	walletID, walletIDErr := uuid.FromString(walletIDAsString)
	if walletIDErr != nil {
		str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
		panic(errors.New(str))
	}

	// retrieve the wallet:
	wal, walErr := walletRepository.RetrieveByID(&walletID)
	if walErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), walErr)
		panic(errors.New(str))
	}

	if !me {
		retUsersPS, retUsersPSErr := userRepository.RetrieveSetByWallet(wal, index, amount)
		if retUsersPSErr != nil {
			str := fmt.Sprintf("there was an error while retrieving a set of users (index: %d, amount: %d, walletID: %s): %s", index, amount, wal.ID().String(), retUsersPSErr.Error())
			return nil, nil, errors.New(str)
		}

		return retUsersPS, nil, nil
	}

	retUser, retUserErr := userRepository.RetrieveByPubKey(conf.WalletPK().PublicKey())
	if retUserErr != nil {
		str := fmt.Sprintf("there was an error while retrieving a user (pubKey: %s, walletID: %s): %s", conf.WalletPK().PublicKey().String(), wal.ID().String(), retUserErr.Error())
		panic(errors.New(str))
	}

	return nil, retUser, nil
}
