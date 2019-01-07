package vote

import (
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

// SDKFunc represents the pledge SDK func
var SDKFunc = struct {
	Create       func(met meta.Meta) *cliapp.Command
	Delete       func() *cliapp.Command
	Retrieve     func() *cliapp.Command
	RetrieveList func() *cliapp.Command
}{
	Create: func(met meta.Meta) *cliapp.Command {
		return create(met)
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
}
