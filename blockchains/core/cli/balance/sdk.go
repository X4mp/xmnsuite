package balance

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the balance SDK func
var SDKFunc = struct {
	Retrieve func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
}
