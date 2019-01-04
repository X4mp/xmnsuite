package information

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Information represents the information instance
type Information interface {
	ID() *uuid.UUID
	GazPricePerKb() int
	ConcensusNeeded() int
	MaxAmountOfValidators() int
}

// Normalized represents the normalized Information instance
type Normalized interface {
}

// Service represents the Information service
type Service interface {
	Save(ins Information) error
}

// Repository represents the Information repository
type Repository interface {
	Retrieve() (Information, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID                    *uuid.UUID
	GazPricePerKb         int
	ConcensusNeeded       int
	MaxAmountOfValidators int
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

// SDKFunc represents the Information SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Information
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateService        func(params CreateServiceParams) Service
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Information {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createInformation(params.ID, params.ConcensusNeeded, params.GazPricePerKb, params.MaxAmountOfValidators)
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
