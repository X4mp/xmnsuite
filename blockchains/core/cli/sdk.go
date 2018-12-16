package cli

import (
	term "github.com/nsf/termbox-go"
	cliapp "github.com/urfave/cli"
)

func reset() {
	term.Sync()
}

// SDKFunc represents the CLI sdk func
var SDKFunc = struct {
	RetrieveGenesis func() *cliapp.Command
}{
	RetrieveGenesis: func() *cliapp.Command {
		return retrieveGenesis()
	},
}
