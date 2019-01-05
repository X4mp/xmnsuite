package information

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the genesis SDK func
var SDKFunc = struct {
	Retrieve func() *cliapp.Command
	Update   func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	Update: func() *cliapp.Command {
		return update()
	},
}
