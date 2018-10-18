package xmn

import (
	"errors"
	"fmt"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * Wallet
 */

type wallet struct {
	WalletID *uuid.UUID `json:"id"`
	CNeeded  string     `json:"concensus_needed"`
}

type jsonWallet struct {
	WalletID string `json:"id"`
	CNeeded  string `json:"concensus_needed"`
}

func createWallet(id *uuid.UUID, concensusNeeded float64) Wallet {
	out := wallet{
		WalletID: id,
		CNeeded:  strconv.FormatFloat(concensusNeeded, 'f', 7, 64),
	}

	return &out
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.WalletID
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() float64 {
	cNeeded, _ := strconv.ParseFloat(app.CNeeded, 64)
	return cNeeded
}

// MarshalJSON marshals the instance to data
func (app *wallet) MarshalJSON() ([]byte, error) {
	return cdc.MarshalJSON(&jsonWallet{
		WalletID: app.WalletID.String(),
		CNeeded:  app.CNeeded,
	})
}

// UnmarshalJSON unmarshals the data to an instance
func (app *wallet) UnmarshalJSON(data []byte) error {
	ptr := new(jsonWallet)
	jsErr := cdc.UnmarshalJSON(data, ptr)
	if jsErr != nil {
		return jsErr
	}

	walletID, walletIDErr := uuid.FromString(ptr.WalletID)
	if walletIDErr != nil {
		return walletIDErr
	}

	app.WalletID = &walletID
	app.CNeeded = ptr.CNeeded
	return nil
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

// SaveAddUserToWalletRequest saves an  add-user-to-wallet request
func (app *walletService) SaveAddUserToWalletRequest(obj AddUserToWalletRequest) error {
	return nil
}

// SaveAddUserToWalletRequestVote saves an  add-user-to-wallet-request vote
func (app *walletService) SaveAddUserToWalletRequestVote(obj AddUserToWalletRequestVote) error {
	return nil
}

// SaveDeleteUserFromWalletRequest saves a  delete-user-from-wallet request
func (app *walletService) SaveDeleteUserFromWalletRequest(obj DelUserFromWalletRequest) error {
	return nil
}

// SaveDeleteUserFromWalletRequestVote saves a  delete-user-from-wallet-request vote
func (app *walletService) SaveDeleteUserFromWalletRequestVote(obj DelUserFromWalletRequestVote) error {
	return nil
}

// SaveDeleteWalletRequest saves a  delete-wallet request
func (app *walletService) SaveDeleteWalletRequest(obj DeleteWalletRequest) error {
	return nil
}

// SaveDeleteWalletRequestVote saves a  delete-wallet-request vote
func (app *walletService) SaveDeleteWalletRequestVote(obj DeleteWalletRequestVote) error {
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

// RetrieveAddUserToWalletRequests retrieves the add-user-to-wallet requests
func (app *walletService) RetrieveAddUserToWalletRequests(index int, amount int) ([]AddUserToWalletRequests, error) {
	return nil, nil
}

// RetrieveAddUserToWalletRequestsByWalletID retrieves an add-user-to-wallet-requests by ID
func (app *walletService) RetrieveAddUserToWalletRequestsByWalletID(id *uuid.UUID) (AddUserToWalletRequests, error) {
	return nil, nil
}

// RetrieveDelUserFromWalletRequests retrieves the delete-user-from-wallet-requests
func (app *walletService) RetrieveDelUserFromWalletRequests(index int, amount int) ([]DelUserFromWalletRequests, error) {
	return nil, nil
}

// RetrieveDelUserFromWalletRequestsByWalletID retrieves a delete-user-from-wallet-requests by walletID
func (app *walletService) RetrieveDelUserFromWalletRequestsByWalletID(id *uuid.UUID) (DelUserFromWalletRequests, error) {
	return nil, nil
}

// RetrieveDelWalletRequests retrieves the delete-user-from-wallet-requests
func (app *walletService) RetrieveDelWalletRequests(index int, amount int) ([]DelWalletRequests, error) {
	return nil, nil
}

// RetrieveDelWalletRequestsByWalletID retrieves a delete-user-from-wallet-requests by walletID
func (app *walletService) RetrieveDelWalletRequestsByWalletID(id *uuid.UUID) (DelWalletRequests, error) {
	return nil, nil
}
