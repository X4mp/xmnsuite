package feature

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

// Feature represents a feature
type Feature interface {
	ID() *uuid.UUID
	Project() project.Project
	Title() string
	Details() string
	CreatedBy() user.User
}

// Normalized represents a normalized feature
type Normalized interface {
}

// Repository represents a feature repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Feature, error)
	RetrieveSetByProject(proj project.Project, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByCreatedByUser(createdBy user.User, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID        *uuid.UUID
	Project   project.Project
	Title     string
	Details   string
	CreatedBy user.User
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Feature SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Feature
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Feature {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createFeature(params.ID, params.Project, params.Title, params.Details, params.CreatedBy)
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
