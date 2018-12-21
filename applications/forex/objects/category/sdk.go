package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

const (
	maxAountOfCharactersForName         = 50
	maxAmountOfCharactersForDescription = 500
)

// Category represents a category
type Category interface {
	ID() *uuid.UUID
	HasParent() bool
	Parent() Category
	Name() string
	Description() string
}

// Repository represents the category repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Category, error)
	RetrieveSetWithNoParent(index int, amount int) (entity.PartialSet, error)
}

// Normalized represents a normalized category
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Parent      Category
	Name        string
	Description string
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

		out, outErr := createCategory(params.ID, params.Name, params.Description)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if curr, ok := ins.(Category); ok {
					out := createStorableCategory(curr)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if cat, ok := ins.(Category); ok {
					keynames := []string{
						retrieveAllCategoriesKeyname(),
					}

					if cat.HasParent() {
						keynames = append(keynames, retrieveCategoriesByParentCategoryIDKeyname(cat.Parent().ID()))
					}

					if !cat.HasParent() {
						keynames = append(keynames, retrieveCcategoriesWithoutParentKeyname())
					}

					return keynames, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
				return nil, errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
}
