package genesis

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/datastore"
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
	CreateRepository     func(ds datastore.DataStore) Repository
	CreateService        func(ds datastore.DataStore) Service
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Genesis {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createGenesis(params.ID, params.GazPricePerKb, params.MaxAmountOfValidators, params.Deposit, params.User)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateRepository: func(ds datastore.DataStore) Repository {
		met := createMetaData()
		entityRepository := entity.SDKFunc.CreateRepository(ds)
		out := createRepository(entityRepository, met)
		return out
	},
	CreateService: func(ds datastore.DataStore) Service {
		met := createMetaData()
		rep := representation()
		entityRepository := entity.SDKFunc.CreateRepository(ds)
		repository := createRepository(entityRepository, met)
		entityService := entity.SDKFunc.CreateService(ds)
		out := createService(entityService, entityRepository, repository, rep)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return representation()
	},
}
