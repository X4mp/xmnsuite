package cryptocurrency

import (
	cliapp "github.com/urfave/cli"
	cryptocurrencycli "github.com/xmnservices/xmnsuite/applications/cryptocurrency/cli"
)

// SDKFunc represents the Cryptocurrency SDK func
var SDKFunc = struct {
	Create func() []cliapp.Command
}{
	Create: func() []cliapp.Command {
		genConf := cryptocurrencycli.SDKFunc.GenerateConfigs()
		spwn := cryptocurrencycli.SDKFunc.Spawn()
		return []cliapp.Command{
			*genConf,
			*spwn,
		}
	},
}
