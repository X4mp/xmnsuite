package validator

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type repository struct {
	entityRepository  entity.Repository
	validatorMetaData entity.MetaData
}

func createRepository(entityRepository entity.Repository, validatorMetaData entity.MetaData) Repository {
	out := repository{
		entityRepository:  entityRepository,
		validatorMetaData: validatorMetaData,
	}

	return &out
}

// RetrieveSet retrieves the validator ordered by their pledge amount
func (app *repository) RetrieveSet(amount int) (entity.PartialSet, error) {
	keyname := retrieveAllValidatorsKeyname()
	return app.entityRepository.RetrieveSetByKeyname(app.validatorMetaData, keyname, 0, amount)
}
