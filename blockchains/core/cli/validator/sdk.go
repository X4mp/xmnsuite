package validator

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the validator SDK func
var SDKFunc = struct {
	Create           func() *cliapp.Command
	Delete           func() *cliapp.Command
	Retrieve         func() *cliapp.Command
	RetrieveList     func() *cliapp.Command
	RetrieveByPledge func() *cliapp.Command
}{
	Create: func() *cliapp.Command {
		return create()
	},
	Delete: func() *cliapp.Command {
		return delete()
	},
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveList: func() *cliapp.Command {
		return retrieveList()
	},
	RetrieveByPledge: func() *cliapp.Command {
		return retrieveByPledge()
	},
}
