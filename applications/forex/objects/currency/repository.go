package currency

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

// RetrieveByID retrieves a currency by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Currency, error) {
	curIns, curInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if curInsErr != nil {
		return nil, curInsErr
	}

	if cur, ok := curIns.(Currency); ok {
		return cur, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a Currency instance", curIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a currency set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllCurrenciesKeyname()
	entityPS, entityPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if entityPSErr != nil {
		return nil, entityPSErr
	}
	return entityPS, nil
}
