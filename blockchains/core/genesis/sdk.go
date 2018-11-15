package genesis

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

// Genesis represents the genesis instance
type Genesis interface {
	ID() *uuid.UUID
	GazPricePerKb() int
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

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	EntityService    entity.Service
	EntityRepository entity.Repository
}

// CreateParams represents the Create params
type CreateParams struct {
	ID                    *uuid.UUID
	GazPricePerKb         int
	MaxAmountOfValidators int
	User                  user.User
	Deposit               deposit.Deposit
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

		out := createGenesis(params.ID, params.GazPricePerKb, params.MaxAmountOfValidators, params.Deposit, params.User)
		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		met := createMetaData()
		out := createRepository(params.EntityRepository, met)
		return out
	},
	CreateService: func(params CreateServiceParams) Service {
		met := createMetaData()
		repository := createRepository(params.EntityRepository, met)
		rep := representation()
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
