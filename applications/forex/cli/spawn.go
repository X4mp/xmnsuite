package cli

import (
	"errors"
	"fmt"
	"time"

	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/applications/forex/commands"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	webserver "github.com/xmnservices/xmnsuite/applications/forex/web"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	coredeposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/configs"
	"github.com/xmnservices/xmnsuite/helpers"
)

func spawn() *cliapp.Command {
	return &cliapp.Command{
		Name:    "spawn",
		Aliases: []string{"s"},
		Usage:   "Spawns the forex blockchain",
		Flags: []cliapp.Flag{
			cliapp.IntFlag{
				Name:  "port",
				Value: 26657,
				Usage: "this is the blockchain port",
			},
			cliapp.IntFlag{
				Name:  "wport",
				Value: 80,
				Usage: "this is the web port",
			},
			cliapp.StringFlag{
				Name:  "dir",
				Value: "./blockchain",
				Usage: "this is the blockchain database path",
			},
			cliapp.StringFlag{
				Name:  "pass",
				Value: "",
				Usage: "this is the password used to decrypt your configuration file",
			},
			cliapp.StringFlag{
				Name:  "file",
				Value: "",
				Usage: "this is the path of your encrypted configuration file",
			},
		},
		Action: func(c *cliapp.Context) error {

			// create the repository:
			repository := configs.SDKFunc.CreateRepository()

			// retrieve the configs:
			retConf, retConfErr := repository.Retrieve(c.String("file"), c.String("pass"))
			if retConfErr != nil {
				return retConfErr
			}

			// spawn the node:
			node := commands.SDKFunc.Spawn(commands.SpawnParams{
				Pass:     c.String("pass"),
				Filename: c.String("file"),
				Dir:      c.String("dir"),
				Port:     c.Int("port"),
			})

			// retrieve the client:
			client, clientErr := node.GetClient()
			if clientErr != nil {
				panic(clientErr)
			}

			// create the entity repository (SDK):
			entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
				PK:          retConf.WalletPK(),
				Client:      client,
				RoutePrefix: "",
			})

			// create the entity service (SDK):
			entityService := entity.SDKFunc.CreateSDKService(entity.CreateSDKServiceParams{
				PK:          retConf.WalletPK(),
				Client:      client,
				RoutePrefix: "",
			})

			// create the account service:
			accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
				PK:          retConf.WalletPK(),
				Client:      client,
				RoutePrefix: "",
			})

			// spawn the web server:
			web := webserver.SDKFunc.Create(webserver.CreateParams{
				Port:           c.Int("wport"),
				EntityService:  entityService,
				AccountService: accountService,
				UserRepository: user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
				BalanceRepository: balance.SDKFunc.CreateRepository(balance.CreateRepositoryParams{
					DepositRepository: coredeposit.SDKFunc.CreateRepository(coredeposit.CreateRepositoryParams{
						EntityRepository: entityRepository,
					}),
					WithdrawalRepository: withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
						EntityRepository: entityRepository,
					}),
				}),
				GenesisRepository: genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
				WalletRepository: wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
				CategoryRepository: category.SDKFunc.CreateRepository(category.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
				CurrencyRepository: currency.SDKFunc.CreateRepository(currency.CreateRepositoryParams{
					EntityRepository: entityRepository,
				}),
			})

			// start the web server:
			err := web.Start()
			if err != nil {
				str := fmt.Sprintf("There was an error while starting the web server: %s", err.Error())
				helpers.Print(str)
			}

			// sleep 1 second before listening to keyboard:
			time.Sleep(time.Second * 1)
			termErr := term.Init()
			if termErr != nil {
				str := fmt.Sprintf("there was an error while enabling the keyboard listening: %s", termErr.Error())
				return errors.New(str)
			}
			defer term.Close()

			// blockchain started, loop until we stop:
			str := fmt.Sprintf("XMN main blockchain spawned, IP: %s\nPress Esc to stop...", client.IP())
			helpers.Print(str)

		keyPressListenerLoop:
			for {
				switch ev := term.PollEvent(); ev.Type {
				case term.EventKey:
					switch ev.Key {
					case term.KeyEsc:
						break keyPressListenerLoop
					}
					break
				}
			}

			// returns:
			return nil
		},
	}
}
