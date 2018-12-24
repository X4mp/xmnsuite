package core

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/datastore"
)

type dependencies struct {
	entityRepository    entity.Repository
	entityService       entity.Service
	groupRepository     group.Repository
	keynameRepository   keyname.Repository
	userRepository      user.Repository
	genesisRepository   genesis.Repository
	genesisService      genesis.Service
	balanceRepository   balance.Repository
	accountService      account.Service
	voteService         active_vote.Service
	developerRepository developer.Repository
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

	accountService := account.SDKFunc.CreateService(account.CreateServiceParams{
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
		EntityService:    entityService,
	})

	voteService := active_vote.SDKFunc.CreateService(active_vote.CreateServiceParams{
		EntityRepository: entityRepository,
		EntityService:    entityService,
	})

	developerRepository := developer.SDKFunc.CreateRepository(ds)

	out := dependencies{
		entityRepository:    entityRepository,
		entityService:       entityService,
		groupRepository:     groupRepository,
		keynameRepository:   keynameRepository,
		userRepository:      userRepository,
		genesisRepository:   genesisRepository,
		genesisService:      genesisService,
		balanceRepository:   balanceRepository,
		accountService:      accountService,
		voteService:         voteService,
		developerRepository: developerRepository,
	}

	return &out
}
