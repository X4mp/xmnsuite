package wallet

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the wallet SDK func
var SDKFunc = struct {
	Create func() *cliapp.Command
	ListMe func() *cliapp.Command
}{
	Create: func() *cliapp.Command {
		return create()
	},
	ListMe: func() *cliapp.Command {
		return listMe()
	},
}
