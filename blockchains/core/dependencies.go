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
	entityRepository := entity.SDKFunc.CreateRepository(entity.CreateRepositoryParams{
		Store: ds,
	})

	entityService := entity.SDKFunc.CreateService(entity.CreateServiceParams{
		Store: ds,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	genesisService := genesis.SDKFunc.CreateService(genesis.CreateServiceParams{
		EntityService:    entityService,
		EntityRepository: entityRepository,
	})

	out := dependencies{
		entityRepository:  entityRepository,
		entityService:     entityService,
		userRepository:    userRepository,
		genesisRepository: genesisRepository,
		genesisService:    genesisService,
	}

	return &out
}
