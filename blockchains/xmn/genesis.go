package xmn

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * Genesis
 */

type genesis struct {
	GzPricePerKb         int            `json:"gaz_price_per_kb"`
	MxAmountOfValidators int            `json:"max_amount_of_validators"`
	Dep                  InitialDeposit `json:"deposit"`
	Tok                  Token          `json:"token"`
}

func createGenesis(gazPricePerKb int, maxAmountOfValidators int, dep InitialDeposit, tok Token) Genesis {
	out := genesis{
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		Dep:                  dep,
		Tok:                  tok,
	}

	return &out
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
func (app *genesis) Deposit() InitialDeposit {
	return app.Dep
}

// Token returns the token
func (app *genesis) Token() Token {
	return app.Tok
}

type storedGenesis struct {
	GzPricePerKb         int `json:"gaz_price_per_kb"`
	MxAmountOfValidators int `json:"max_amount_of_validators"`
}

func createStoredGenesis(gen Genesis) *storedGenesis {
	out := storedGenesis{
		GzPricePerKb:         gen.GazPricePerKb(),
		MxAmountOfValidators: gen.MaxAmountOfValidators(),
	}

	return &out
}

type genesisService struct {
	keyname           string
	store             datastore.DataStore
	walService        WalletService
	tokService        TokenService
	initialDepService InitialDepositService
}

func createGenesisService(
	store datastore.DataStore,
	walService WalletService,
	tokService TokenService,
	initialDepService InitialDepositService,
) GenesisService {
	out := genesisService{
		keyname:           "genesis",
		store:             store,
		walService:        walService,
		tokService:        tokService,
		initialDepService: initialDepService,
	}

	return &out
}

// Save save the genesis instance
func (app *genesisService) Save(obj Genesis) error {
	// make sure the instance does not exists already:
	_, retErr := app.Retrieve()
	if retErr == nil {
		return errors.New("the Genesis instance already exists")
	}

	// save the token:
	saveTokErr := app.tokService.Save(obj.Token())
	if saveTokErr != nil {
		str := fmt.Sprintf("there was an error while saving the Token instance, in the Genesis instance: %s", saveTokErr.Error())
		return errors.New(str)
	}

	// save the initial deposit:
	saveInitialDepErr := app.initialDepService.Save(obj.Deposit())
	if saveInitialDepErr != nil {
		str := fmt.Sprintf("there was an error while saving the InitialDeposit instance, in the Genesis instance: %s", saveInitialDepErr.Error())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keyname,
		Obj: createStoredGenesis(obj),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Genesis instance")
	}

	return nil
}

// Retrieve retrieves the genesis instance
func (app *genesisService) Retrieve() (Genesis, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keyname,
		Obj: new(storedGenesis),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		return nil, errors.New("there was an error while retrieving the Genesis instance")
	}

	// cast the instance:
	if storedGen, ok := obj.Obj.(*storedGenesis); ok {
		// retrieve the token:
		tok, tokErr := app.tokService.Retrieve()
		if tokErr != nil {
			str := fmt.Sprintf("there was an error while retrieving the Token in the Genesis instance: %s", tokErr.Error())
			return nil, errors.New(str)
		}

		// retrieve the initial deposit:
		initialDep, initialDepErr := app.initialDepService.Retrieve()
		if initialDepErr != nil {
			str := fmt.Sprintf("there was an error while retrieving the InitialDeposit in the Genesis instance: %s", initialDepErr.Error())
			return nil, errors.New(str)
		}

		out := createGenesis(storedGen.GzPricePerKb, storedGen.MxAmountOfValidators, initialDep, tok)
		return out, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a Genesis instance")
}
