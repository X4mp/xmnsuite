package category

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

// Category represents a category
type Category interface {
	ID() *uuid.UUID
	Title() string
	Description() string
	HasParent() bool
	Parent() Category
}

// Normalized represents a normalized category
type Normalized interface {
}

// Repository represents the category repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Category, error)
	RetrieveSetByParent(parent Category, index int, amount int) (entity.PartialSet, error)
	RetrieveSetWithoutParent(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Title       string
	Description string
	Parent      Category
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Category SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Category
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Category {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		if params.Parent != nil {
			out, outErr := createCategory(params.ID, params.Title, params.Description)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := createCategoryWithParent(params.ID, params.Title, params.Description, params.Parent)
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
