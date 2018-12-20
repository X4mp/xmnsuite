package blockchains

import (
	"encoding/hex"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/commands"
	"github.com/xmnservices/xmnsuite/configs"
)

type blockchain struct {
	startCommand commands.Command
}

func createBlockchain(
	port int,
	name string,
	namespace string,
	id *uuid.UUID,
	conf configs.Configs,
	blockchainRootDirectory string,
	databaseFilePath string,
	peers []string,
	genTrs genesis.Genesis,
	met meta.Meta,
) (Blockchain, error) {
	configParams := commands.CreateParams{
		Constants: commands.CreateConstantsParams{
			Namespace:     namespace,
			Name:          name,
			ID:            id.String(),
			RoutePrefix:   "",
			RouterRoleKey: "router-role-key",
		},
		Port:                    port,
		NodePrivateKey:          hex.EncodeToString(conf.NodePK().Bytes()),
		BlockchainRootDirectory: blockchainRootDirectory,
		DatabaseFilePath:        databaseFilePath,
		Meta:                    met,
	}

	// if the directory does not exists, execute the genesis command:
	if _, err := os.Stat(blockchainRootDirectory); os.IsNotExist(err) {
		genCmd := commands.SDKFunc.CreateGenesis(commands.CreateGenesisParams{
			Configs:            configParams,
			RootPrivateKey:     conf.WalletPK().String(),
			GenesisTransaction: genTrs,
		})

		// execute the command:
		genNode, genNodeErr := genCmd.Execute()
		if genNodeErr != nil {
			return nil, genNodeErr
		}

		// stop the node:
		genNodeStopErr := genNode.Stop()
		if genNodeStopErr != nil {
			return nil, genNodeStopErr
		}

		// execute the command:
		genNode, genNodeErr = genCmd.Execute()
		if genNodeErr != nil {
			return nil, genNodeErr
		}

		// stop the node:
		genNodeStopErr = genNode.Stop()
		if genNodeStopErr != nil {
			return nil, genNodeStopErr
		}
	}

	// execute the start command:
	startCmd := commands.SDKFunc.CreateStart(commands.CreateStartParams{
		Configs: configParams,
		Peers:   peers,
	})

	out := blockchain{
		startCommand: startCmd,
	}

	return &out, nil
}

// Start starts the blockchain
func (app *blockchain) Start() (applications.Node, error) {
	// execute the command:
	node, nodeErr := app.startCommand.Execute()
	if nodeErr != nil {
		return nil, nodeErr
	}

	return node, nil
}
