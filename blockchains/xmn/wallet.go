package xmn

import (
	"errors"
	"fmt"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * Wallet
 */

type wallet struct {
	WalletID *uuid.UUID `json:"id"`
	CNeeded  string     `json:"concensus_needed"`
	cNeeded  float64
	Usrs     []User `json:"users"`
}

type jsonWallet struct {
	WalletID string `json:"id"`
	CNeeded  string `json:"concensus_needed"`
	Usrs     []User `json:"users"`
}

func createWallet(id *uuid.UUID, concensusNeeded float64, usrs []User) Wallet {
	out := wallet{
		WalletID: id,
		CNeeded:  strconv.FormatFloat(concensusNeeded, 'f', 7, 64),
		cNeeded:  concensusNeeded,
		Usrs:     usrs,
	}

	return &out
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.WalletID
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() float64 {
	return app.cNeeded
}

// Users returns the users
func (app *wallet) Users() []User {
	return app.Usrs
}

// MarshalJSON marshals the instance to data
func (app *wallet) MarshalJSON() ([]byte, error) {
	return cdc.MarshalJSON(&jsonWallet{
		WalletID: app.WalletID.String(),
		CNeeded:  app.CNeeded,
		Usrs:     app.Usrs,
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

	cNeeded, cNeededErr := strconv.ParseFloat(ptr.CNeeded, 64)
	if cNeededErr != nil {
		return cNeededErr
	}

	app.WalletID = &walletID
	app.CNeeded = ptr.CNeeded
	app.cNeeded = cNeeded
	app.Usrs = ptr.Usrs
	return nil
}

/*
 * WalletService
 */

type walletService struct {
	store datastore.DataStore
}

func createWalletService(store datastore.DataStore) WalletService {
	out := walletService{
		store: store,
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
		Key: fmt.Sprintf("wallet:by_id:%s", wallet.ID().String()),
		Obj: wallet,
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Wallet instance")
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

// RetrieveByID retrieves a wallet by ID
func (app *walletService) RetrieveByID(id *uuid.UUID) (Wallet, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: fmt.Sprintf("wallet:by_id:%s", id.String()),
		Obj: new(wallet),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		str := fmt.Sprintf("there was an error while retrieving the Wallet instance (ID: %s)", id.String())
		return nil, errors.New(str)
	}

	// cast the instance:
	if wal, ok := obj.Obj.(*wallet); ok {
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

/*
 * User
 */

type user struct {
	PKey  string `json:"pubkey"`
	Shres int    `json:"shares"`
}

func createUser(pubKey crypto.PublicKey, shares int) User {
	out := user{
		PKey:  pubKey.String(),
		Shres: shares,
	}

	return &out
}

// PubKey returns the PublicKey
func (app *user) PubKey() crypto.PublicKey {
	return crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: app.PKey,
	})
}

// Shares returns the shares
func (app *user) Shares() int {
	return app.Shres
}
