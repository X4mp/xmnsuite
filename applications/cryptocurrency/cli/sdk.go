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
	GenerateConfigs func() *cliapp.Command
	Spawn           func() *cliapp.Command
}{
	GenerateConfigs: func() *cliapp.Command {
		return generateConfigs()
	},
	Spawn: func() *cliapp.Command {
		return spawn()
	},
}
