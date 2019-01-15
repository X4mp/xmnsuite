package project

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

// Project represents a project
type Project interface {
	ID() *uuid.UUID
	Proposal() proposal.Proposal
}

// Normalized represents the normalized project
type Normalized interface {
}

// Repository represents the project repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Project, error)
	RetrieveByProposal(prop proposal.Proposal) (Project, error)
	RetrieveSetByCategory(cat category.Category, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID       *uuid.UUID
	Proposal proposal.Proposal
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Project SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Project
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Project {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createProject(
			params.ID,
			params.Proposal,
		)

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
