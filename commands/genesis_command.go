package commands

import (
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/datastore"
)

type genesisCommand struct {
	conf GenesisConfigs
}

func createGenesisCommand(conf GenesisConfigs) (Command, error) {
	out := genesisCommand{
		conf: conf,
	}

	return &out, nil
}

// Execute executes a genesis command
func (app *genesisCommand) Execute() (applications.Node, error) {
	// create the blockchain:
	conf := app.conf.Configs()
	cons := conf.Constants()
	blkchain := tendermint.SDKFunc.CreateBlockchain(tendermint.CreateBlockchainParams{
		Namespace: cons.Namespace(),
		Name:      cons.Name(),
		ID:        cons.ID(),
		PrivKey:   conf.NodePrivateKey(),
	})

	// create the blockchain service:
	service := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: conf.BlockchainRootDirectory(),
	})

	// save the blockchain:
	saveErr := service.Save(blkchain)
	if saveErr != nil {
		return nil, saveErr
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
		RootPubKey:    app.conf.RootPrivateKey().PublicKey(),
		Store:         store,
		Meta:          conf.Meta(),
	})

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(conf.Port(), nil, conf.BlockchainRootDirectory(), blkchain, apps)
	if nodeErr != nil {
		return nil, nodeErr
	}

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		return nil, startNodeErr
	}

	// get the client:
	client, clientErr := node.GetClient()
	if clientErr != nil {
		return nil, clientErr
	}

	// create the genesis service:
	entityService := entity.SDKFunc.CreateSDKService(entity.CreateSDKServiceParams{
		PK:          app.conf.RootPrivateKey(),
		Client:      client,
		RoutePrefix: cons.RoutePrefix(),
	})

	genesisService := genesis.SDKFunc.CreateService(genesis.CreateServiceParams{
		EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
			PK:          app.conf.RootPrivateKey(),
			Client:      client,
			RoutePrefix: cons.RoutePrefix(),
		}),
		EntityService: entityService,
	})

	// save the genesis transaction:
	saveGenErr := genesisService.Save(app.conf.GenesisTransaction())
	if saveGenErr != nil {
		return nil, saveGenErr
	}

	// save the request group and keyname:
	entReqs := app.conf.Configs().Meta().WriteOnEntityRequest()
	for _, entReq := range entReqs {
		grp := group.SDKFunc.Create(group.CreateParams{
			Name: entReq.RequestedBy().MetaData().Keyname(),
		})

		mp := entReq.Map()
		keynameRepresentation := keyname.SDKFunc.CreateRepresentation()
		for _, oneRepresentation := range mp {
			kname := keyname.SDKFunc.Create(keyname.CreateParams{
				Name:  oneRepresentation.MetaData().Keyname(),
				Group: grp,
			})

			// save the keyname:
			saveKeynameErr := entityService.Save(kname, keynameRepresentation)
			if saveKeynameErr != nil {
				return nil, saveKeynameErr
			}
		}
	}

	// everything worked, return:
	return node, nil
}
