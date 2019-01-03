package cli

import (
	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/wallet"
)

func reset() {
	term.Sync()
}

// SDKFunc represents the CLI sdk func
var SDKFunc = struct {
	RetrieveGenesis func() *cliapp.Command
	Wallet          func() *cliapp.Command
}{
	RetrieveGenesis: func() *cliapp.Command {
		return retrieveGenesis()
	},
	Wallet: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "This is the group of commands to work with wallets",
			Subcommands: []cliapp.Command{
				*wallet.SDKFunc.Create(),
			},
		}
	},
}
