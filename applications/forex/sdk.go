package forex

import (
	cliapp "github.com/urfave/cli"
	forexcli "github.com/xmnservices/xmnsuite/applications/forex/cli"
)

// SDKFunc represents the forex SDK func
var SDKFunc = struct {
	Create func() []cliapp.Command
}{
	Create: func() []cliapp.Command {
		genConf := forexcli.SDKFunc.GenerateConfigs()
		spwn := forexcli.SDKFunc.Spawn()
		return []cliapp.Command{
			*genConf,
			*spwn,
		}
	},
}
