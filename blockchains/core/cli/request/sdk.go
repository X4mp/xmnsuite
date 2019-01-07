package request

import (
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/cli/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Active  func(met meta.Meta) *cliapp.Command
	Group   func() *cliapp.Command
	Keyname func() *cliapp.Command
}{
	Active: func(met meta.Meta) *cliapp.Command {
		return &cliapp.Command{
			Name:    "active",
			Aliases: []string{"a"},
			Usage:   "This is the group of commands to work with active requests",
			Subcommands: []cliapp.Command{
				*active.SDKFunc.Retrieve(),
				*active.SDKFunc.RetrieveList(),
				*active.SDKFunc.RetrieveListFromWallet(),
				*active.SDKFunc.Vote(met),
			},
		}
	},
	Group: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "This is the group of commands to work with request groups",
			Subcommands: []cliapp.Command{
				*group.SDKFunc.Retrieve(),
				*group.SDKFunc.RetrieveList(),
			},
		}
	},
	Keyname: func() *cliapp.Command {
		return &cliapp.Command{
			Name:    "keyname",
			Aliases: []string{"k"},
			Usage:   "This is the group of commands to work with request keynames",
			Subcommands: []cliapp.Command{
				*keyname.SDKFunc.Retrieve(),
				*keyname.SDKFunc.RetrieveList(),
			},
		}
	},
}
