package active

import (
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/request/active/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Retrieve               func() *cliapp.Command
	RetrieveList           func() *cliapp.Command
	RetrieveListFromWallet func() *cliapp.Command
	Vote                   func(met meta.Meta) *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveList: func() *cliapp.Command {
		return retrieveList()
	},
	RetrieveListFromWallet: func() *cliapp.Command {
		return retrieveListFromWallet()
	},
	Vote: func(met meta.Meta) *cliapp.Command {
		return &cliapp.Command{
			Name:    "vote",
			Aliases: []string{"v"},
			Usage:   "This is the group of commands to work with votes on requests",
			Subcommands: []cliapp.Command{
				*vote.SDKFunc.Create(met),
				*vote.SDKFunc.Delete(),
				*vote.SDKFunc.Retrieve(),
				*vote.SDKFunc.RetrieveList(),
			},
		}
	},
}
