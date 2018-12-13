package commands

import (
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/datastore"
)

type startCommand struct {
	conf StartConfigs
}

func createStartCommand(conf StartConfigs) (Command, error) {
	out := startCommand{
		conf: conf,
	}

	return &out, nil
}

// Execute executes the command
func (app *startCommand) Execute() (applications.Node, error) {
	// create the blockchain path:
	conf := app.conf.Configs()
	cons := conf.Constants()
	blkchainPath := tendermint.SDKFunc.CreatePath(tendermint.CreatePathParams{
		Namespace: cons.Namespace(),
		Name:      cons.Name(),
		ID:        cons.ID(),
	})

	// create the blockchain service:
	service := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: conf.BlockchainRootDirectory(),
	})

	// retrieve the blockchain:
	blkchain, blkchainErr := service.Retrieve(blkchainPath)
	if blkchainErr != nil {
		panic(blkchainErr)
	}

	// create the datastore:
	store := datastore.SDKFunc.CreateStoredDataStore(datastore.StoredDataStoreParams{
		FilePath: conf.DatabaseFilePath(),
	})

	// create the core applications:
	apps := core.SDKFunc.Create(core.CreateParams{
		Namespace:     cons.Namespace(),
		Name:          cons.Name(),
		ID:            cons.ID(),
		Port:          conf.Port(),
		NodePK:        conf.NodePrivateKey(),
		RootDir:       conf.BlockchainRootDirectory(),
		RoutePrefix:   cons.RoutePrefix(),
		RouterRoleKey: cons.RouterRoleKey(),
		Store:         store,
		Meta:          conf.Meta(),
	})

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// create the peers if any:
	seeds := []string{}
	if app.conf.HasPeers() {
		peers := app.conf.Peers()
		for _, onePeer := range peers {
			seeds = append(seeds, onePeer.String())
		}
	}

	// spawn the node:
	node, nodeErr := appService.Spawn(conf.Port(), seeds, conf.BlockchainRootDirectory(), blkchain, apps)
	if nodeErr != nil {
		return nil, nodeErr
	}

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		return nil, startNodeErr
	}

	// everything worked, return:
	return node, nil
}
