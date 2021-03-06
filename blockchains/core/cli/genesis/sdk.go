package genesis

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the genesis SDK func
var SDKFunc = struct {
	Retrieve func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
}
