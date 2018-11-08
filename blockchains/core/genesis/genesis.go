package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
)

/*
 * Genesis
 */

type genesis struct {
	UUID                 *uuid.UUID      `json:"id"`
	GzPricePerKb         int             `json:"gaz_price_per_kb"`
	MxAmountOfValidators int             `json:"max_amount_of_validators"`
	Dep                  deposit.Deposit `json:"deposit"`
}

func createGenesis(id *uuid.UUID, gazPricePerKb int, maxAmountOfValidators int, dep deposit.Deposit) Genesis {
	out := genesis{
		UUID:                 id,
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		Dep:                  dep,
	}

	return &out
}

func createGenesisFromNormalized(ins *normalizedGenesis) (Genesis, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	depIns, depInsErr := deposit.SDKFunc.CreateMetaData().Denormalize()(ins.Deposit)
	if depInsErr != nil {
		return nil, depInsErr
	}

	if dep, ok := depIns.(deposit.Deposit); ok {
		out := createGenesis(&id, ins.GzPricePerKb, ins.MxAmountOfValidators, dep)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (app *genesis) ID() *uuid.UUID {
	return app.UUID
}

// GazPricePerKb returns the gazPricePerKb
func (app *genesis) GazPricePerKb() int {
	return app.GzPricePerKb
}

// MaxAmountOfValidators returns the maxAmountOfValidators
func (app *genesis) MaxAmountOfValidators() int {
	return app.MxAmountOfValidators
}

// Deposit returns the initial deposit
func (app *genesis) Deposit() deposit.Deposit {
	return app.Dep
}
