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
	catIns, catInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if catInsErr != nil {
		return nil, catInsErr
	}

	if cat, ok := catIns.(Category); ok {
		return cat, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", catIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetWithNoParent retrieves a category set
func (app *repository) RetrieveSetWithNoParent(index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllCategoriesKeyname(),
		retrieveCcategoriesWithoutParentKeyname(),
	}

	catPS, catPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if catPSErr != nil {
		return nil, catPSErr
	}

	return catPS, nil
}
