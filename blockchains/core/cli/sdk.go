package cli

import (
	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/information"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/wallet"
)

func reset() {
	term.Sync()
}

// SDKFunc represents the CLI sdk func
var SDKFunc = struct {
	Spawn       func() *cliapp.Command
	Genesis     func() *cliapp.Command
	Information func() *cliapp.Command
	Wallet      func() *cliapp.Command
}{
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
			Aliases: []string{"g"},
			Usage:   "This is the group of commands to work with the information instance",
			Subcommands: []cliapp.Command{
				*information.SDKFunc.Retrieve(),
				*information.SDKFunc.Update(),
			},
		}
	},
	Wallet: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "This is the group of commands to work with wallets",
			Subcommands: []cliapp.Command{
				*wallet.SDKFunc.Create(),
				*wallet.SDKFunc.ListMe(),
			},
		}
	},
}
