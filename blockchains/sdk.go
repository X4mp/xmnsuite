package blockchains

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/configs"
)

// Blockchain represents the blockchain application
type Blockchain interface {
	Start() (applications.Node, error)
}

// CreateParams represents the create params
type CreateParams struct {
	Port                    int
	Name                    string
	Namespace               string
	ID                      string
	Conf                    configs.Configs
	BlockchainRootDirectory string
	DatabaseFilePath        string
	Peers                   []string
	GenesisTransaction      genesis.Genesis
	Meta                    meta.Meta
}

// SDKFunc represents the blockchains SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Blockchain
}{
	Create: func(params CreateParams) Blockchain {
		id, idErr := uuid.FromString(params.ID)
		if idErr != nil {
			panic(idErr)
		}

		out, outErr := createBlockchain(
			params.Port,
			params.Name,
			params.Namespace,
			&id,
			params.Conf,
			params.BlockchainRootDirectory,
			params.DatabaseFilePath,
			params.Peers,
			params.GenesisTransaction,
			params.Meta,
		)

		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
