package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type wallet struct {
	WalletID *uuid.UUID `json:"id"`
	CNeeded  int        `json:"concensus_needed"`
}

type jsonWallet struct {
	WalletID string `json:"id"`
	CNeeded  int    `json:"concensus_needed"`
}

func createWallet(id *uuid.UUID, concensusNeeded int) Wallet {
	out := wallet{
		WalletID: id,
		CNeeded:  concensusNeeded,
	}

	return &out
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.WalletID
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() int {
	return app.CNeeded
}

type walletPartialSet struct {
	Wals  []Wallet `json:"wallets"`
	Idx   int      `json:"index"`
	TotAm int      `json:"total_amount"`
}

func createWalletPartialSet(wals []Wallet, idx int, totAm int) WalletPartialSet {
	out := walletPartialSet{
		Wals:  wals,
		Idx:   idx,
		TotAm: totAm,
	}

	return &out
}

// Wallets returns the wallets
func (obj *walletPartialSet) Wallets() []Wallet {
	return obj.Wals
}

// Index returns the index
func (obj *walletPartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *walletPartialSet) Amount() int {
	return len(obj.Wals)
}

// TotalAmount returns the total amount
func (obj *walletPartialSet) TotalAmount() int {
	return obj.TotAm
}

/*
 * WalletService
 */

type walletService struct {
	keyname string
	store   datastore.DataStore
}

func createWalletService(store datastore.DataStore) WalletService {
	out := walletService{
		keyname: "wallets",
		store:   store,
	}

	return &out
}

// Save saves a wallet
func (app *walletService) Save(wallet Wallet) error {
	// make sure the instance does not exists already:
	_, retErr := app.RetrieveByID(wallet.ID())
	if retErr == nil {
		str := fmt.Sprintf("the Wallet instance (ID: %s) already exists", wallet.ID().String())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: fmt.Sprintf("%s:by_id:%s", app.keyname, wallet.ID().String()),
		Obj: wallet,
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Wallet instance")
	}

	// add the wallet to the list:
	amountAdded := app.store.Sets().Add(app.keyname, wallet.ID().String())
	if amountAdded != 1 {
		str := fmt.Sprintf("there was an error while saving the Wallet ID (%s) to the wallets set (wallets)", wallet.ID().String())
		return errors.New(str)
	}

	return nil
}

// Retrieve retrieves a list of wallets
func (app *walletService) Retrieve(index int, amount int) (WalletPartialSet, error) {
	// retrieve wallet uuids:
	ids := app.store.Sets().Retrieve(app.keyname, index, amount)

	// retrieve the wallets:
	wallets := []Wallet{}
	for _, oneIDAsString := range ids {
		// cast the ID:
		id, idErr := uuid.FromString(oneIDAsString.(string))
		if idErr != nil {
			str := fmt.Sprintf("there is an element (%s) in the wallets set (keyname: wallets) that is not a valid wallet UUID", oneIDAsString)
			return nil, errors.New(str)
		}

		wal, walErr := app.RetrieveByID(&id)
		if walErr != nil {
			return nil, walErr
		}

		wallets = append(wallets, wal)
	}

	// retrieve the total amount of elements in the keyname:
	totalAmount := app.store.Sets().Len(app.keyname)

	// return:
	out := createWalletPartialSet(wallets, index, totalAmount)
	return out, nil
}

// RetrieveByID retrieves a wallet by ID
func (app *walletService) RetrieveByID(id *uuid.UUID) (Wallet, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: fmt.Sprintf("%s:by_id:%s", app.keyname, id.String()),
		Obj: new(wallet),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		str := fmt.Sprintf("there was an error while retrieving the Wallet instance (ID: %s)", id.String())
		return nil, errors.New(str)
	}

	// cast the instance:
	if wal, ok := obj.Obj.(Wallet); ok {
		return wal, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a Wallet instance")
}
