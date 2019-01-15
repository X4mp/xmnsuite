package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	metaData         entity.MetaData
	entityRepository entity.Repository
}

func createRepository(metaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		metaData:         metaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves a category by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Category, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if cat, ok := ins.(Category); ok {
		return cat, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByParent retrieves a category partial set by parent category
func (app *repository) RetrieveSetByParent(parent Category, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllCategoriesKeyname(),
		retrieveCategoryByParentCategoryKeyname(parent),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}

// RetrieveSetWithoutParent retrieves a category partial set without parent
func (app *repository) RetrieveSetWithoutParent(index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllCategoriesKeyname(),
		retrieveCategoryWithoutParentKeyname(),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}
