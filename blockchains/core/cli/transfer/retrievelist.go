package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/configs"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrievelist() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve-list",
		Aliases: []string{"rl"},
		Usage:   "Retrieves a list of transfers related to a wallet",
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
				Usage: "This is the walletid we want to retrieve transfers from. If not set, all transfers are retrieved. (Optional)",
			},
			cliapp.BoolFlag{
				Name:  "incoming",
				Usage: "If set to true, only show the incoming transfers",
			},
			cliapp.BoolFlag{
				Name:  "outgoing",
				Usage: "If set to true, only show the outgoing transfers",
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

			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			transferRepository := transfer.SDKFunc.CreateRepository(transfer.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			withdrawalRepository := withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			// retrieve the genesis:
			gen, genErr := genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr)
				panic(errors.New(str))
			}

			// set the index and amount:
			trsfPS, trsfs, trsfsErr := retrieveRightTransfers(c, conf, gen, transferRepository, walletRepository, depositRepository, withdrawalRepository)
			if trsfsErr != nil {
				panic(trsfsErr)
			}

			if trsfPS != nil {
				helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
					Ins: trsfPS,
				})

				return nil
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: trsfs,
			})

			// returns:
			return nil
		},
	}
}

func retrieveRightTransfers(
	c *cliapp.Context,
	conf configs.Configs,
	gen genesis.Genesis,
	transferRepository transfer.Repository,
	walletRepository wallet.Repository,
	depositRepository deposit.Repository,
	withdrawalRepository withdrawal.Repository,
) (entity.PartialSet, []transfer.Transfer, error) {
	// get the variables:
	index := c.Int("index")
	amount := c.Int("amount")
	walletIDAsString := c.String("walletid")
	isIncoming := c.Bool("incoming")
	isOutgoing := c.Bool("outgoing")

	// if there is no wallet:
	if walletIDAsString == "" {
		retTrsfPS, retTrsfPSErr := transferRepository.RetrieveSet(index, amount)
		if retTrsfPSErr != nil {
			str := fmt.Sprintf("there was an error while retrieving a set of transfers (index: %d, amount: %d): %s", index, amount, retTrsfPSErr.Error())
			return nil, nil, errors.New(str)
		}

		return retTrsfPS, nil, nil
	}

	// parse the walletID:
	walletID, walletIDErr := uuid.FromString(walletIDAsString)
	if walletIDErr != nil {
		str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
		return nil, nil, errors.New(str)
	}

	// retrieve the wallet:
	wal, walErr := walletRepository.RetrieveByID(&walletID)
	if walErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), walErr)
		return nil, nil, errors.New(str)
	}

	// get the token:
	tok := gen.Deposit().Token()

	// if the transaction is incoming:
	if isIncoming {
		// retrieve the deposits related to our wallet:
		deps, depsErr := depositRepository.RetrieveSetByToWalletAndToken(wal, tok)
		if depsErr != nil {
			str := fmt.Sprintf("there was a problem while retrieving a deposits (walletID: %s, tokenID: %s): %s", wal.ID().String(), tok.ID().String(), depsErr.Error())
			return nil, nil, errors.New(str)
		}

		// retrieve the transfers related to the deposits:
		trsfs := []transfer.Transfer{}
		for _, oneDep := range deps {
			trsf, trsfErr := transferRepository.RetrieveByDeposit(oneDep)
			if trsfErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the transfer (DepositID: %s): %s", oneDep.ID().String(), trsfErr.Error())
				return nil, nil, errors.New(str)
			}

			trsfs = append(trsfs, trsf)
		}

		return nil, trsfs, nil
	}

	if isOutgoing {
		// retrieve the withdrawals related to our wallet:
		withs, withsErr := withdrawalRepository.RetrieveSetByFromWalletAndToken(wal, tok)
		if withsErr != nil {
			str := fmt.Sprintf("there was a problem while retrieving a withdrawals (walletID: %s, tokenID: %s): %s", wal.ID().String(), tok.ID().String(), withsErr.Error())
			return nil, nil, errors.New(str)
		}

		// retrieve the transfers related to the deposits:
		trsfs := []transfer.Transfer{}
		for _, oneWith := range withs {
			trsf, trsfErr := transferRepository.RetrieveByWithdrawal(oneWith)
			if trsfErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the transfer (WithdrawalIDS: %s): %s", oneWith.ID().String(), trsfErr.Error())
				return nil, nil, errors.New(str)
			}

			trsfs = append(trsfs, trsf)
		}

		return nil, trsfs, nil
	}

	trsfPS, trsfPSErr := transferRepository.RetrieveSet(index, amount)
	if trsfPSErr != nil {
		str := fmt.Sprintf("there was an error while retrieving transfer partial set (index: %d, amount: %d): %s", index, amount, trsfPSErr.Error())
		return nil, nil, errors.New(str)
	}

	return trsfPS, nil, nil
}
