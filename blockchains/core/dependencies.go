package core

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/datastore"
)

type dependencies struct {
	entityRepository    entity.Repository
	entityService       entity.Service
	groupRepository     group.Repository
	keynameRepository   keyname.Repository
	walletRepository    wallet.Repository
	userRepository      user.Repository
	genesisRepository   genesis.Repository
	genesisService      genesis.Service
	balanceRepository   balance.Repository
	voteService         active_vote.Service
	affiliateRepository affiliates.Repository
	validatorRepository validator.Repository
}

func createDependencies(ds datastore.DataStore) *dependencies {
	entityRepository := entity.SDKFunc.CreateRepository(ds)
	entityService := entity.SDKFunc.CreateService(ds)

	groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	genesisService := genesis.SDKFunc.CreateService(genesis.CreateServiceParams{
		EntityRepository: entityRepository,
		EntityService:    entityService,
	})

	balanceRepository := balance.SDKFunc.CreateRepository(balance.CreateRepositoryParams{
		Datastore: ds,
	})

	walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	voteService := active_vote.SDKFunc.CreateService(active_vote.CreateServiceParams{
		EntityRepository: entityRepository,
		EntityService:    entityService,
	})

	affiliateRepository := affiliates.SDKFunc.CreateRepository(affiliates.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	validatorRepository := validator.SDKFunc.CreateRepository(validator.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	out := dependencies{
		entityRepository:    entityRepository,
		entityService:       entityService,
		groupRepository:     groupRepository,
		keynameRepository:   keynameRepository,
		walletRepository:    walletRepository,
		userRepository:      userRepository,
		genesisRepository:   genesisRepository,
		genesisService:      genesisService,
		balanceRepository:   balanceRepository,
		voteService:         voteService,
		affiliateRepository: affiliateRepository,
		validatorRepository: validatorRepository,
	}

	return &out
}
