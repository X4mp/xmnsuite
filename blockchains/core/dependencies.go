package core

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/retrievers/balance"
	"github.com/xmnservices/xmnsuite/datastore"
)

type dependencies struct {
	entityRepository  entity.Repository
	entityService     entity.Service
	userRepository    user.Repository
	genesisRepository genesis.Repository
	genesisService    genesis.Service
	balanceRepository balance.Repository
}

func createDependencies(ds datastore.DataStore) *dependencies {
	entityRepository := entity.SDKFunc.CreateRepository(ds)
	entityService := entity.SDKFunc.CreateService(ds)
	userRepository := user.SDKFunc.CreateRepository(ds)
	genesisRepository := genesis.SDKFunc.CreateRepository(ds)
	genesisService := genesis.SDKFunc.CreateService(ds)
	balanceRepository := balance.SDKFunc.CreateRepository(ds)

	out := dependencies{
		entityRepository:  entityRepository,
		entityService:     entityService,
		userRepository:    userRepository,
		genesisRepository: genesisRepository,
		genesisService:    genesisService,
		balanceRepository: balanceRepository,
	}

	return &out
}
