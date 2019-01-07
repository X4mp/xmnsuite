package user

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the genesis SDK func
var SDKFunc = struct {
	Retrieve     func() *cliapp.Command
	RetrieveList func() *cliapp.Command
	Save         func() *cliapp.Command
	SaveToWallet func() *cliapp.Command
	Delete       func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveList: func() *cliapp.Command {
		return retrievelist()
	},
	Save: func() *cliapp.Command {
		return save()
	},
	SaveToWallet: func() *cliapp.Command {
		return saveToWallet()
	},
	Delete: func() *cliapp.Command {
		return delete()
	},
}
