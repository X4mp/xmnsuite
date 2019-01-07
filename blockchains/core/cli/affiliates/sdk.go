package affiliates

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Retrieve         func() *cliapp.Command
	RetrieveByWallet func() *cliapp.Command
	RetrieveList     func() *cliapp.Command
	Save             func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveByWallet: func() *cliapp.Command {
		return retrieveByWallet()
	},
	RetrieveList: func() *cliapp.Command {
		return retrieveList()
	},
	Save: func() *cliapp.Command {
		return save()
	},
}
