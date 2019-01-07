package pledge

import (
	cliapp "github.com/urfave/cli"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Retrieve     func() *cliapp.Command
	RetrieveFrom func() *cliapp.Command
	RetrieveTo   func() *cliapp.Command
	Create       func() *cliapp.Command
	Delete       func() *cliapp.Command
}{
	Retrieve: func() *cliapp.Command {
		return retrieve()
	},
	RetrieveFrom: func() *cliapp.Command {
		return retrieveFrom()
	},
	RetrieveTo: func() *cliapp.Command {
		return retrieveTo()
	},
	Create: func() *cliapp.Command {
		return create()
	},
	Delete: func() *cliapp.Command {
		return delete()
	},
}
