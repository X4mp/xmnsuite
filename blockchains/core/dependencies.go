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
	genesisRepository genesis.Repository
	userRepository    user.Repository
}

func createDependencies(ds datastore.DataStore) *dependencies {
	entityRepository := entity.SDKFunc.CreateRepository(entity.CreateRepositoryParams{
		Store: ds,
	})

	entityService := entity.SDKFunc.CreateService(entity.CreateServiceParams{
		Store: ds,
	})

	genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	out := dependencies{
		entityRepository:  entityRepository,
		entityService:     entityService,
		genesisRepository: genesisRepository,
		userRepository:    userRepository,
	}

	return &out
}
