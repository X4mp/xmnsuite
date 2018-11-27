package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
)

/*
 * Genesis
 */

type genesis struct {
	UUID                 *uuid.UUID      `json:"id"`
	ConNeeded            int             `json:"concensus_needed"`
	DevConNeeded         int             `json:"developer_concensus_needed"`
	GzPricePerKb         int             `json:"gaz_price_per_kb"`
	MxAmountOfValidators int             `json:"max_amount_of_validators"`
	Usr                  user.User       `json:"user"`
	Dep                  deposit.Deposit `json:"deposit"`
}

func createGenesis(id *uuid.UUID, concensusNeeded int, devConcensusNeeded int, gazPricePerKb int, maxAmountOfValidators int, dep deposit.Deposit, usr user.User) (Genesis, error) {

	out := genesis{
		UUID:                 id,
		ConNeeded:            concensusNeeded,
		DevConNeeded:         devConcensusNeeded,
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		Usr:                  usr,
		Dep:                  dep,
	}

	return &out, nil
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

	usrIns, usrInsErr := user.SDKFunc.CreateMetaData().Denormalize()(ins.User)
	if usrInsErr != nil {
		return nil, usrInsErr
	}

	if dep, ok := depIns.(deposit.Deposit); ok {
		if usr, ok := usrIns.(user.User); ok {
			return createGenesis(&id, ins.ConcensusNeeded, ins.DeveloperConcensusNeeded, ins.GzPricePerKb, ins.MxAmountOfValidators, dep, usr)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", usrIns.ID().String())
		return nil, errors.New(str)
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

// ConcensusNeeded returns the concensusNeeded
func (app *genesis) ConcensusNeeded() int {
	return app.ConNeeded
}

// DeveloperConcensusNeeded returns the developer concensusNeeded
func (app *genesis) DeveloperConcensusNeeded() int {
	return app.DevConNeeded
}

// MaxAmountOfValidators returns the maxAmountOfValidators
func (app *genesis) MaxAmountOfValidators() int {
	return app.MxAmountOfValidators
}

// User returns the user
func (app *genesis) User() user.User {
	return app.Usr
}

// Deposit returns the initial deposit
func (app *genesis) Deposit() deposit.Deposit {
	return app.Dep
}
