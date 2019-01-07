package group

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Retrieve     func() *cliapp.Command
	RetrieveList func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveList: func() *cliapp.Command {
		return retrieveList()
	},
}
