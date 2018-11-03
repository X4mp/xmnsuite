package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type storedTransfer struct {
	ID      string `json:"id"`
	Amount  int    `json:"amount"`
	Content string `json:"content"`
	PubKey  string `json:"public_key"`
}

func createStoredTransfer(trx Transfer) *storedTransfer {
	out := storedTransfer{
		ID:      trx.ID().String(),
		Amount:  trx.Amount(),
		Content: trx.Content(),
		PubKey:  trx.PubKey().String(),
	}

	return &out
}

type transfer struct {
	UUID *uuid.UUID       `json:"id"`
	Am   int              `json:"amount"`
	Cnt  string           `json:"content"`
	PKey crypto.PublicKey `json:"public_key"`
}

func createTransfer(id *uuid.UUID, amount int, content string, pubKey crypto.PublicKey) Transfer {
	out := transfer{
		UUID: id,
		Am:   amount,
		Cnt:  content,
		PKey: pubKey,
	}

	return &out
}

// ID returns the ID
func (obj *transfer) ID() *uuid.UUID {
	return obj.UUID
}

// Amount returns the amount
func (obj *transfer) Amount() int {
	return obj.Am
}

// Content returns the content
func (obj *transfer) Content() string {
	return obj.Cnt
}

// PubKey returns the public key
func (obj *transfer) PubKey() crypto.PublicKey {
	return obj.PKey
}

type transferPartialSet struct {
	Trx   []Transfer `json:"transfers"`
	Idx   int        `json:"index"`
	TotAm int        `json:"total_amount"`
}

func createTransferPartialSet(trx []Transfer, index int, totalAmount int) TransferPartialSet {
	out := transferPartialSet{
		Trx:   trx,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out
}

// Transfers returns the transfers
func (obj *transferPartialSet) Transfers() []Transfer {
	return obj.Trx
}

// Index returns the index
func (obj *transferPartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *transferPartialSet) Amount() int {
	return len(obj.Trx)
}

// TotalAmount returns the totalAmount
func (obj *transferPartialSet) TotalAmount() int {
	return obj.TotAm
}

type transferService struct {
	keyname string
	store   datastore.DataStore
}

func createTransferService(store datastore.DataStore) TransferService {
	out := transferService{
		keyname: "transfers",
		store:   store,
	}

	return &out
}

// Save saves a Transfer instance
func (app *transferService) Save(trx Transfer) error {
	// make sure the instance does not exists already:
	_, retErr := app.RetrieveByID(trx.ID())
	if retErr == nil {
		str := fmt.Sprintf("the Transfer instance (ID: %s) already exists", trx.ID().String())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keynameByID(trx.ID()),
		Obj: createStoredTransfer(trx),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Transfer instance")
	}

	// add the pubKey to the sets:
	keynameByPubKey := app.keynameByPublicKey(trx.PubKey())
	amountAdded := app.store.Sets().Add(keynameByPubKey, trx.ID().String())
	if amountAdded != 1 {
		// revert:
		app.store.Sets().Del(keynameByPubKey, trx.ID().String())

		str := fmt.Sprintf("there was an error while saving the Transfer ID (%s) to the transfer sets... reverting", trx.ID().String())
		return errors.New(str)
	}

	return nil
}

// Retrieve retrieves TransferPartialSet instance
func (app *transferService) Retrieve(index int, amount int) (TransferPartialSet, error) {
	return nil, nil
}

// RetrieveByID retrieves a Transfer instance by ID
func (app *transferService) RetrieveByID(id *uuid.UUID) (Transfer, error) {
	return nil, nil
}

// RetrieveByPublicKey retrieves a Transfer instance by public key
func (app *transferService) RetrieveByPublicKey(pubKey crypto.PublicKey) (Transfer, error) {
	return nil, nil
}

func (app *transferService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *transferService) keynameByPublicKey(pubKey crypto.PublicKey) string {
	return fmt.Sprintf("%s:by_pubkey:%s", app.keyname, pubKey.String())
}
