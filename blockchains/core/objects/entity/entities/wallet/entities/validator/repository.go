package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
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

// RetrieveByID retrieves a validator instance by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Validator, error) {
	valIns, valInsErr := app.entityRepository.RetrieveByID(app.validatorMetaData, id)
	if valInsErr != nil {
		return nil, valInsErr
	}

	if val, ok := valIns.(Validator); ok {
		return val, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", valIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a validator set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllValidatorsKeyname()
	return app.entityRepository.RetrieveSetByKeyname(app.validatorMetaData, keyname, index, amount)
}

// RetrieveByPledge retrieves a validator by pledge
func (app *repository) RetrieveByPledge(pldge pledge.Pledge) (Validator, error) {
	keynames := []string{
		retrieveAllValidatorsKeyname(),
		retrieveValidatorsByPledgeKeyname(pldge),
	}

	valIns, valInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.validatorMetaData, keynames)
	if valInsErr != nil {
		return nil, valInsErr
	}

	if val, ok := valIns.(Validator); ok {
		return val, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", valIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetOrderedByPledgeAmount retrieves the validators ordered by pledge amount
func (app *repository) RetrieveSetOrderedByPledgeAmount(index int, amount int) ([]Validator, error) {
	valsPS, valsPSErr := app.RetrieveSet(0, -1)
	if valsPSErr != nil {
		return nil, valsPSErr
	}

	vals := []Validator{}
	valIns := valsPS.Instances()
	for _, oneValIns := range valIns {
		if val, ok := oneValIns.(Validator); ok {
			vals = append(vals, val)
			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", oneValIns.ID().String())
		return nil, errors.New(str)
	}

	return orderValPSByPledge(vals, index, amount)
}
