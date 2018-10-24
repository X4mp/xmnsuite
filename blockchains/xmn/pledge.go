package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type storedPledge struct {
	ID   string `json:"id"`
	From string `json:"from_wallet_id"`
	To   string `json:"to_wallet_id"`
	Am   int    `json:"amount"`
}

func createStoredPledge(pledge Pledge) *storedPledge {
	out := storedPledge{
		ID:   pledge.ID().String(),
		From: pledge.From().ID().String(),
		To:   pledge.To().ID().String(),
		Am:   pledge.Amount(),
	}

	return &out
}

type pledge struct {
	UUID       *uuid.UUID `json:"id"`
	FromWallet Wallet     `json:"from"`
	ToWallet   Wallet     `json:"to"`
	Am         int        `json:"amount"`
}

func createPledge(id *uuid.UUID, frm Wallet, to Wallet, amount int) Pledge {
	out := pledge{
		UUID:       id,
		FromWallet: frm,
		ToWallet:   to,
		Am:         amount,
	}

	return &out
}

// ID returns the ID
func (obj *pledge) ID() *uuid.UUID {
	return obj.UUID
}

// From returns the from wallet
func (obj *pledge) From() Wallet {
	return obj.FromWallet
}

// To returns the to wallet
func (obj *pledge) To() Wallet {
	return obj.ToWallet
}

// Amount returns the amount
func (obj *pledge) Amount() int {
	return obj.Am
}

type pledgePartialSet struct {
	Plds  []Pledge `json:"pledges"`
	Idx   int      `json:"index"`
	TotAm int      `json:"total_amount"`
}

func createPledgePartialSet(pledges []Pledge, idx int, totAm int) PledgePartialSet {
	out := pledgePartialSet{
		Plds:  pledges,
		Idx:   idx,
		TotAm: totAm,
	}

	return &out
}

// Pledges returns the pledges
func (obj *pledgePartialSet) Pledges() []Pledge {
	return obj.Plds
}

// Index returns the index
func (obj *pledgePartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *pledgePartialSet) Amount() int {
	return len(obj.Plds)
}

// TotalAmount returns the totalAmount
func (obj *pledgePartialSet) TotalAmount() int {
	return obj.TotAm
}

type pledgeService struct {
	keyname       string
	store         datastore.DataStore
	walletService WalletService
}

func createPledgeService(store datastore.DataStore, walletService WalletService) PledgeService {
	out := pledgeService{
		keyname:       "pledges",
		store:         store,
		walletService: walletService,
	}

	return &out
}

// Save saves a Pledge instance
func (app *pledgeService) Save(pledge Pledge) error {
	// make sure the instance does not exists already:
	_, retErr := app.RetrieveByID(pledge.ID())
	if retErr == nil {
		str := fmt.Sprintf("the Pledge instance (ID: %s) already exists", pledge.ID().String())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keynameByID(pledge.ID()),
		Obj: createStoredPledge(pledge),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Pledge instance")
	}

	// create the set keys:
	keys := []string{
		app.keynameByFromWalletID(pledge.From().ID()),
		app.keynameByToWalletID(pledge.To().ID()),
	}

	// add the wallet to the sets:
	amountAdded := app.store.Sets().AddMul(keys, pledge.ID().String())
	if amountAdded != 1 {
		// revert:
		app.store.Sets().DelMul(keys, pledge.ID().String())

		str := fmt.Sprintf("there was an error while saving the Pledge ID (%s) to the pledges sets... reverting", pledge.ID().String())
		return errors.New(str)
	}

	return nil
}

// RetrieveByID retrieves a Pledge by ID
func (app *pledgeService) RetrieveByID(id *uuid.UUID) (Pledge, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keynameByID(id),
		Obj: new(storedPledge),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		return nil, errors.New("there was an error while retrieving the Pledge instance")
	}

	// cast the instance:
	if storedPledge, ok := obj.Obj.(*storedPledge); ok {
		return app.FromStoredToPledge(storedPledge)
	}

	return nil, errors.New("the retrieved data cannot be casted to a Pledge instance")
}

// RetrieveByFromWalletID retrieves a PledgePartialSet instance by the from WalletID
func (app *pledgeService) RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error) {
	keyname := app.keynameByFromWalletID(fromWalletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

// RetrieveByToWalletID retrieves a PledgePartialSet instance by the to WalletID
func (app *pledgeService) RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error) {
	keyname := app.keynameByToWalletID(toWalletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

func (app *pledgeService) FromStoredToPledge(stored *storedPledge) (Pledge, error) {
	// cast the ID:
	id, idErr := uuid.FromString(stored.ID)
	if idErr != nil {
		return nil, idErr
	}

	// cast the fromWalletID:
	fromWalletID, fromWalletIDErr := uuid.FromString(stored.From)
	if fromWalletIDErr != nil {
		return nil, fromWalletIDErr
	}

	// cast the toWalletID:
	toWalletID, toWalletIDErr := uuid.FromString(stored.To)
	if toWalletIDErr != nil {
		return nil, toWalletIDErr
	}

	// retrieve the from wallet:
	fromWallet, fromWalletErr := app.walletService.RetrieveByID(&fromWalletID)
	if fromWalletErr != nil {
		return nil, fromWalletErr
	}

	// retrieve the to wallet:
	toWallet, toWalletErr := app.walletService.RetrieveByID(&toWalletID)
	if toWalletErr != nil {
		return nil, toWalletErr
	}

	out := createPledge(&id, fromWallet, toWallet, stored.Am)
	return out, nil
}

func (app *pledgeService) retrievePartialSetByKeyname(keyname string, index int, amount int) (PledgePartialSet, error) {
	// retrieve pledge uuids:
	ids := app.store.Sets().Retrieve(keyname, index, amount)

	// retrieve the pledges:
	pledges := []Pledge{}
	for _, oneIDAsString := range ids {
		// cast the ID:
		id, idErr := uuid.FromString(oneIDAsString.(string))
		if idErr != nil {
			str := fmt.Sprintf("there is an element (%s) in the pledge set (keyname: %s) that is not a valid Pledge UUID", keyname, oneIDAsString)
			return nil, errors.New(str)
		}

		pledge, pledgeErr := app.RetrieveByID(&id)
		if pledgeErr != nil {
			return nil, pledgeErr
		}

		pledges = append(pledges, pledge)
	}

	// retrieve the total amount of elements in the keyname:
	totalAmount := app.store.Sets().Len(keyname)

	// return:
	out := createPledgePartialSet(pledges, index, totalAmount)
	return out, nil
}

func (app *pledgeService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *pledgeService) keynameByFromWalletID(fromWalletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_from_wallet_id:%s", app.keyname, fromWalletID.String())
}

func (app *pledgeService) keynameByToWalletID(toWalletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_to_wallet_id:%s", app.keyname, toWalletID.String())
}
