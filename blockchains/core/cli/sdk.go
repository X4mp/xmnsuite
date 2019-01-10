package cli

import (
	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/affiliates"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/information"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

func reset() {
	term.Sync()
}

// SDKFunc represents the CLI sdk func
var SDKFunc = struct {
	Config 		func() *cliapp.Command
	Spawn       func() *cliapp.Command
	Genesis     func() *cliapp.Command
	Information func() *cliapp.Command
	Balance     func() *cliapp.Command
	User        func() *cliapp.Command
	Transfer    func() *cliapp.Command
	Pledge      func() *cliapp.Command
	Validator   func() *cliapp.Command
	Wallet      func() *cliapp.Command
	Affiliates  func() *cliapp.Command
	Request     func(met meta.Meta) *cliapp.Command
}{
	Config: func() *cliapp.Command {
		return generateConfig()
	},
	Spawn: func() *cliapp.Command {
		return spawn()
	},
	Genesis: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "genesis",
			Aliases: []string{"g"},
			Usage:   "This is the group of commands to work with the genesis instance",
			Subcommands: []cliapp.Command{
				*genesis.SDKFunc.Retrieve(),
			},
		}
	},
	Information: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "information",
			Aliases: []string{"i"},
			Usage:   "This is the group of commands to work with the information instance",
			Subcommands: []cliapp.Command{
				*information.SDKFunc.Retrieve(),
				*information.SDKFunc.Update(),
			},
		}
	},
	Balance: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "balance",
			Aliases: []string{"b"},
			Usage:   "This is the group of commands to work with the balance of wallets",
			Subcommands: []cliapp.Command{
				*balance.SDKFunc.Retrieve(),
			},
		}
	},
	User: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "This is the group of commands to work with users",
			Subcommands: []cliapp.Command{
				*user.SDKFunc.Retrieve(),
				*user.SDKFunc.RetrieveList(),
				*user.SDKFunc.Save(),
				*user.SDKFunc.SaveToWallet(),
				*user.SDKFunc.Delete(),
			},
		}
	},
	Transfer: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "transfer",
			Aliases: []string{"t"},
			Usage:   "This is the group of commands to work with transfers",
			Subcommands: []cliapp.Command{
				*transfer.SDKFunc.Retrieve(),
				*transfer.SDKFunc.RetrieveList(),
				*transfer.SDKFunc.Save(),
			},
		}
	},
	Pledge: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "pledge",
			Aliases: []string{"p"},
			Usage:   "This is the group of commands to work with pledges",
			Subcommands: []cliapp.Command{
				*pledge.SDKFunc.Retrieve(),
				*pledge.SDKFunc.RetrieveFrom(),
				*pledge.SDKFunc.RetrieveTo(),
				*pledge.SDKFunc.Create(),
				*pledge.SDKFunc.Delete(),
			},
		}
	},
	Validator: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "validator",
			Aliases: []string{"v"},
			Usage:   "This is the group of commands to work with validators",
			Subcommands: []cliapp.Command{
				*validator.SDKFunc.Retrieve(),
				*validator.SDKFunc.RetrieveByPledge(),
				*validator.SDKFunc.RetrieveList(),
				*validator.SDKFunc.Create(),
				*validator.SDKFunc.Delete(),
			},
		}
	},
	Wallet: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "This is the group of commands to work with wallets",
			Subcommands: []cliapp.Command{
				*wallet.SDKFunc.Retrieve(),
				*wallet.SDKFunc.RetrieveList(),
			},
		}
	},
	Affiliates: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "affiliates",
			Aliases: []string{"w"},
			Usage:   "This is the group of commands to work with affiliates",
			Subcommands: []cliapp.Command{
				*affiliates.SDKFunc.Retrieve(),
				*affiliates.SDKFunc.RetrieveByWallet(),
				*affiliates.SDKFunc.RetrieveList(),
				*affiliates.SDKFunc.Save(),
			},
		}
	},
	Request: func(met meta.Meta) *cliapp.Command {
		return &cliapp.Command{
			Name:    "request",
			Aliases: []string{"r"},
			Usage:   "This is the group of commands to work with requests",
			Subcommands: []cliapp.Command{
				*request.SDKFunc.Active(met),
				*request.SDKFunc.Group(),
				*request.SDKFunc.Keyname(),
			},
		}
	},
}
