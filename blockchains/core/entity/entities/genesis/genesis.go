package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
)

/*
 * Genesis
 */

type genesis struct {
	UUID                 *uuid.UUID      `json:"id"`
	GzPricePerKb         int             `json:"gaz_price_per_kb"`
	MxAmountOfValidators int             `json:"max_amount_of_validators"`
	Usr                  user.User       `json:"user"`
	Dep                  deposit.Deposit `json:"deposit"`
}

func createGenesis(id *uuid.UUID, gazPricePerKb int, maxAmountOfValidators int, dep deposit.Deposit, usr user.User) (Genesis, error) {

	// the user must have enough shares in order to fill the concensus, on genesis:
	if usr.Shares() < dep.To().ConcensusNeeded() {
		str := fmt.Sprintf("the genesis user (ID: %s) does not have enough shares (%d) in order to convert the concensus of the genesis wallet (ID: %s ConcensusNeeded: %d)", usr.ID().String(), usr.Shares(), dep.To().ID(), dep.To().ConcensusNeeded())
		return nil, errors.New(str)
	}

	out := genesis{
		UUID:                 id,
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
			return createGenesis(&id, ins.GzPricePerKb, ins.MxAmountOfValidators, dep, usr)
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
