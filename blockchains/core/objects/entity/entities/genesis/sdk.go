package genesis

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Genesis represents the genesis instance
type Genesis interface {
	ID() *uuid.UUID
	GazPricePerKb() int
	ConcensusNeeded() int
	MaxAmountOfValidators() int
	User() user.User
	Deposit() deposit.Deposit
}

// Normalized represents the normalized Genesis instance
type Normalized interface {
}

// Service represents the Genesis service
type Service interface {
	Save(ins Genesis) error
}

// Repository represents the Genesis repository
type Repository interface {
	Retrieve() (Genesis, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID                    *uuid.UUID
	GazPricePerKb         int
	ConcensusNeeded       int
	MaxAmountOfValidators int
	User                  user.User
	Deposit               deposit.Deposit
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	Datastore        datastore.DataStore
	EntityRepository entity.Repository
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	Datastore        datastore.DataStore
	EntityRepository entity.Repository
	EntityService    entity.Service
}

// SDKFunc represents the Genesis SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Genesis
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateService        func(params CreateServiceParams) Service
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Genesis {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createGenesis(params.ID, params.ConcensusNeeded, params.GazPricePerKb, params.MaxAmountOfValidators, params.Deposit, params.User)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		if params.Datastore != nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.Datastore)
		}

		met := createMetaData()
		out := createRepository(params.EntityRepository, met)
		return out
	},
	CreateService: func(params CreateServiceParams) Service {
		if params.Datastore != nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.Datastore)
			params.EntityService = entity.SDKFunc.CreateService(params.Datastore)
		}

		met := createMetaData()
		rep := representation()
		repository := createRepository(params.EntityRepository, met)
		out := createService(params.EntityService, params.EntityRepository, repository, rep)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return representation()
	},
}
