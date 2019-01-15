package proposal

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

// Proposal represents a proposal
type Proposal interface {
	ID() *uuid.UUID
	Title() string
	Description() string
	Details() string
	Category() category.Category
	ManagerPledgeNeeded() int
	LinkerPledgeNeeded() int
}

// Normalized represents a normalized proposal
type Normalized interface {
}

// Repository represents the proposal repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Proposal, error)
	RetrieveSetByCategory(cat category.Category, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID                  *uuid.UUID
	Title               string
	Description         string
	Details             string
	Category            category.Category
	ManagerPledgeNeeded int
	LinkerPledgeNeeded  int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Proposal SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Proposal
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Proposal {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createProposal(params.ID, params.Title, params.Description, params.Details, params.Category, params.ManagerPledgeNeeded, params.LinkerPledgeNeeded)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return representation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
}
