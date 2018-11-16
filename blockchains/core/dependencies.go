package core

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
	"github.com/xmnservices/xmnsuite/datastore"
)

type dependencies struct {
	entityRepository  entity.Repository
	entityService     entity.Service
	userRepository    user.Repository
	genesisRepository genesis.Repository
	genesisService    genesis.Service
}

func createDependencies(ds datastore.DataStore) *dependencies {
	entityRepository := entity.SDKFunc.CreateRepository(ds)
	entityService := entity.SDKFunc.CreateService(ds)
	userRepository := user.SDKFunc.CreateRepository(ds)
	genesisRepository := genesis.SDKFunc.CreateRepository(ds)
	genesisService := genesis.SDKFunc.CreateService(ds)

	out := dependencies{
		entityRepository:  entityRepository,
		entityService:     entityService,
		userRepository:    userRepository,
		genesisRepository: genesisRepository,
		genesisService:    genesisService,
	}

	return &out
}
