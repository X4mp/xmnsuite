package commands

import (
	"net"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
)

// SpawnMainParams represents the spawn main params
type SpawnMainParams struct {
	Pass     string
	Filename string
	Dir      string
	Port     int
}

// RetrieveGenesisParams retrieve the genesis transaction params
type RetrieveGenesisParams struct {
	Pass     string
	Filename string
	IP       net.IP
	Port     int
}

// SDKFunc represents the commands SDK func
var SDKFunc = struct {
	RetrieveGenesis func(params RetrieveGenesisParams) genesis.Genesis
}{
	RetrieveGenesis: func(params RetrieveGenesisParams) genesis.Genesis {
		out, outErr := retrieveGenesis(params.Pass, params.Filename, params.IP, params.Port)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
