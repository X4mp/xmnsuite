package balance

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/helpers"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	core_helpers "github.com/xmnservices/xmnsuite/helpers"
)

func retrieve() *cliapp.Command {
	return &cliapp.Command{
		Name:    "retrieve",
		Aliases: []string{"r"},
		Usage:   "Retrieve the balance of a wallet",
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
				Usage: "This is the walletid we want to retrieve the balance from",
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

			genRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: entityRepository,
			})

			balanceRepository := balance.SDKFunc.CreateRepository(balance.CreateRepositoryParams{
				DepositRepository: deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
				WithdrawalRepository: withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
			})

			// retrieve the genesis:
			gen, genErr := genRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
				panic(errors.New(str))
			}

			// parse the walletID:
			walletIDAsString := c.String("walletid")
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

			// retrieve the balance:
			tok := gen.Deposit().Token()
			bal, balErr := balanceRepository.RetrieveByWalletAndToken(wal, tok)
			if balErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the balance (WalletID: %s, TokenID: %s): %s", wal.ID().String(), tok.ID().String(), balErr.Error())
				panic(errors.New(str))
			}

			helpers.SDKFunc.PrintSuccessWithInstance(helpers.PrintSuccessWithInstanceParams{
				Ins: bal,
			})

			// returns:
			return nil
		},
	}
}
