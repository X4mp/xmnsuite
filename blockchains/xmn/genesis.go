package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * Genesis
 */

type genesis struct {
	GzPricePerKb         int            `json:"gaz_price_per_kb"`
	MxAmountOfValidators int            `json:"max_amount_of_validators"`
	Devs                 Wallet         `json:"developers"`
	Dep                  InitialDeposit `json:"deposit"`
	Tok                  Token          `json:"token"`
}

func createGenesis(gazPricePerKb int, maxAmountOfValidators int, devs Wallet, dep InitialDeposit, tok Token) Genesis {
	out := genesis{
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		Devs:                 devs,
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

// Developers returns the developers wallet
func (app *genesis) Developers() Wallet {
	return app.Devs
}

// Deposit returns the initial deposit
func (app *genesis) Deposit() InitialDeposit {
	return app.Dep
}

// Token returns the token
func (app *genesis) Token() Token {
	return app.Tok
}

/*
 * Token
 */

type token struct {
	Sym  string `json:"symbol"`
	Nme  string `json:"name"`
	Desc string `json:"description"`
}

func createToken(symbol string, name string, desc string) Token {
	out := token{
		Sym:  symbol,
		Nme:  name,
		Desc: desc,
	}

	return &out
}

// Symbol returns the symbol
func (app *token) Symbol() string {
	return app.Sym
}

// Name returns the name
func (app *token) Name() string {
	return app.Nme
}

// Description returns the description
func (app *token) Description() string {
	return app.Desc
}

/*
 * InitialDeposit
 */

type initialDeposit struct {
	ToWallet Wallet `json:"to"`
	Am       int    `json:"amount"`
}

func createInitialDeposit(wallet Wallet, amount int) InitialDeposit {
	out := initialDeposit{
		ToWallet: wallet,
		Am:       amount,
	}

	return &out
}

// To returns the to wallet
func (app *initialDeposit) To() Wallet {
	return app.ToWallet
}

// Amount returns the amount
func (app *initialDeposit) Amount() int {
	return app.Am
}

/*
 * GenesisService
 */

type storedGenesis struct {
	GzPricePerKb         int            `json:"gaz_price_per_kb"`
	MxAmountOfValidators int            `json:"max_amount_of_validators"`
	DevsWalletID         *uuid.UUID     `json:"developers_wallet_id"`
	Dep                  InitialDeposit `json:"deposit"`
	Tok                  Token          `json:"token"`
}

func createStoredGenesis(gen Genesis) *storedGenesis {
	out := storedGenesis{
		GzPricePerKb:         gen.GazPricePerKb(),
		MxAmountOfValidators: gen.MaxAmountOfValidators(),
		DevsWalletID:         gen.Developers().ID(),
		Dep:                  gen.Deposit(),
		Tok:                  gen.Token(),
	}

	return &out
}

type genesisService struct {
	keyname    string
	store      datastore.DataStore
	walService WalletService
}

func createGenesisService(store datastore.DataStore, walService WalletService) GenesisService {
	out := genesisService{
		keyname:    "genesis-instance",
		store:      store,
		walService: walService,
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

	// save the wallet:
	saveWalletErr := app.walService.Save(obj.Developers())
	if saveWalletErr != nil {
		str := fmt.Sprintf("there was an error while saving the developer's wallet: %s", saveWalletErr.Error())
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
		// retrieve the wallet:
		wal, wallErr := app.walService.RetrieveByID(storedGen.DevsWalletID)
		if wallErr != nil {
			str := fmt.Sprintf("the developer's wallet (ID: %s) cannot be retrieved: %s", storedGen.DevsWalletID.String(), wallErr.Error())
			return nil, errors.New(str)
		}

		out := createGenesis(storedGen.GzPricePerKb, storedGen.MxAmountOfValidators, wal, storedGen.Dep, storedGen.Tok)
		return out, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a Genesis instance")
}
