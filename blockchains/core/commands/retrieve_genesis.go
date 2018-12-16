package commands

import (
	"net"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/configs"
)

func retrieveGenesis(pass string, filename string, ip net.IP, port int) (genesis.Genesis, error) {
	// create the repository:
	repository := configs.SDKFunc.CreateRepository()

	// retrieve the configs:
	retConf, retConfErr := repository.Retrieve(filename, pass)
	if retConfErr != nil {
		return nil, retConfErr
	}

	// connect to the blockchain:
	client := tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
		IP:   ip,
		Port: port,
	})

	// create the genesis repository:
	genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
		EntityRepository: entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
			PK:          retConf.WalletPK(),
			Client:      client,
			RoutePrefix: "",
		}),
	})

	// retrieve the genesis:
	gen, genErr := genesisRepository.Retrieve()
	if genErr != nil {
		return nil, genErr
	}

	return gen, nil
}
