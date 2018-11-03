package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type storedTransferRequest struct {
	FromUserID string `json:"from_user_id"`
	Amount     string `json:"amount"`
	PubKey     string `json:"public_key"`
	Reason     string `json:"reason"`
}

func createStoredTransferRequest(req TransferRequest) *storedTransferRequest {
	out := storedTransferRequest{
		FromUserID: req.From().ID().String(),
		Amount:     req.Amount(),
		PubKey:     req.PubKey().String(),
		Reason:     req.Reason(),
	}

	return &out
}

type transferRequest struct {
	Frm   User             `json:"from"`
	Am    string           `json:"amount"`
	PbKey crypto.PublicKey `json:"public_key"`
	Rs    string           `json:"reason"`
}

func createTransferRequest(from User, amount string, pubKey crypto.PublicKey, reason string) TransferRequest {
	out := transferRequest{
		Frm:   from,
		Am:    amount,
		PbKey: pubKey,
		Rs:    reason,
	}

	return &out
}

// From returns the from user
func (obj *transferRequest) From() User {
	return obj.Frm
}

// Amount returns the ciphered amount
func (obj *transferRequest) Amount() string {
	return obj.Am
}

// PubKey returns the PubKey used to encrypt the amount
func (obj *transferRequest) PubKey() crypto.PublicKey {
	return obj.PbKey
}

// Reason returns the transfer reason
func (obj *transferRequest) Reason() string {
	return obj.Rs
}

type storedSignedTransferRequest struct {
	ID  string                 `json:"id"`
	Req *storedTransferRequest `json:"transfer_request"`
	Sig string                 `json:"signature"`
}

func createStoredSignedTransferRequest(transfer SignedTransferRequest) *storedSignedTransferRequest {
	out := storedSignedTransferRequest{
		ID:  transfer.ID().String(),
		Req: createStoredTransferRequest(transfer.Request()),
		Sig: transfer.Signature().String(),
	}

	return &out
}

type signedTransferRequest struct {
	UUID *uuid.UUID           `json:"id"`
	Req  TransferRequest      `json:"transfer_request"`
	Sig  crypto.RingSignature `json:"signature"`
}

func createSignedTransferRequest(id *uuid.UUID, req TransferRequest, sig crypto.RingSignature) SignedTransferRequest {
	out := signedTransferRequest{
		UUID: id,
		Req:  req,
		Sig:  sig,
	}

	return &out
}

// ID returns the ID
func (obj *signedTransferRequest) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the TransferRequest
func (obj *signedTransferRequest) Request() TransferRequest {
	return obj.Req
}

// Signature returns the RingSignature
func (obj *signedTransferRequest) Signature() crypto.RingSignature {
	return obj.Sig
}

type signedTransferPartialSet struct {
	Reqs  []SignedTransferRequest `json:"transfer_requests"`
	Idx   int                     `json:"index"`
	TotAm int                     `json:"total_amount"`
}

func createSignedTransferPartialSet(reqs []SignedTransferRequest, index int, totalAmount int) SignedTransferRequestPartialSet {
	out := signedTransferPartialSet{
		Reqs:  reqs,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out
}

// Requests returns the transfer requests
func (obj *signedTransferPartialSet) Requests() []SignedTransferRequest {
	return obj.Reqs
}

// Index returns the index
func (obj *signedTransferPartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *signedTransferPartialSet) Amount() int {
	return len(obj.Reqs)
}

// TotalAmount returns the total amount
func (obj *signedTransferPartialSet) TotalAmount() int {
	return obj.TotAm
}

type signedTransferRequestService struct {
	keyname     string
	store       datastore.DataStore
	userService UserService
}

func createSignedTransferRequestService(store datastore.DataStore, userService UserService) SignedTransferRequestService {
	out := signedTransferRequestService{
		keyname:     "signed-transfers",
		store:       store,
		userService: userService,
	}

	return &out
}

// Save saves a SignedTransferRequest instance
func (app *signedTransferRequestService) Save(req SignedTransferRequest) error {
	// make sure the instance does not exists already:
	_, retErr := app.RetrieveByID(req.ID())
	if retErr == nil {
		str := fmt.Sprintf("the TransferRequest instance (ID: %s) already exists", req.ID().String())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keynameByID(req.ID()),
		Obj: createStoredSignedTransferRequest(req),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the TransferRequest instance")
	}

	// add the wallet to the sets:
	keynameByWalletID := app.keynameByFromWalletID(req.Request().From().Wallet().ID())
	amountAdded := app.store.Sets().Add(keynameByWalletID, req.ID().String())
	if amountAdded != 1 {
		// revert:
		app.store.Sets().Del(keynameByWalletID, req.ID().String())

		str := fmt.Sprintf("there was an error while saving the TransferRequest ID (%s) to the transferrequest sets... reverting", req.ID().String())
		return errors.New(str)
	}

	return nil
}

// RetrieveByID retrieves a SignedTransferRequest by ID
func (app *signedTransferRequestService) RetrieveByID(id *uuid.UUID) (SignedTransferRequest, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keynameByID(id),
		Obj: new(storedTransferRequest),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		return nil, errors.New("there was an error while retrieving the TransferRequest instance")
	}

	// cast the instance:
	if storedReq, ok := obj.Obj.(*storedSignedTransferRequest); ok {
		return app.FromStoredToSignedTransferRequest(storedReq)
	}

	return nil, errors.New("the retrieved data cannot be casted to a TransferRequest instance")
}

// RetrieveByFromWalletID retrieves the SignedTransferRequestPartialSet sent from the given walletID
func (app *signedTransferRequestService) RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (SignedTransferRequestPartialSet, error) {
	keyname := app.keynameByID(fromWalletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

// FromStoredToSignedTransferRequest converts a StoredSignedTransferRequest to a SignedTransferRequest
func (app *signedTransferRequestService) FromStoredToSignedTransferRequest(stored *storedSignedTransferRequest) (SignedTransferRequest, error) {
	// cast the ID:
	id, idErr := uuid.FromString(stored.ID)
	if idErr != nil {
		return nil, idErr
	}

	// cast the ring signature:
	ringSig := crypto.SDKFunc.CreateRingSig(crypto.CreateRingSigParams{
		RingSigAsString: stored.Sig,
	})

	// get the request:
	storedTrxReq := stored.Req

	// cast the fromUserID:
	fromUserID, fromUserIDErr := uuid.FromString(storedTrxReq.FromUserID)
	if fromUserIDErr != nil {
		return nil, fromUserIDErr
	}

	// retrieve the user:
	fromUser, fromUserErr := app.userService.RetrieveByID(&fromUserID)
	if fromUserErr != nil {
		return nil, fromUserErr
	}

	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: storedTrxReq.PubKey,
	})

	trxReq := createTransferRequest(fromUser, storedTrxReq.Amount, pubKey, storedTrxReq.Reason)
	out := createSignedTransferRequest(&id, trxReq, ringSig)
	return out, nil
}

func (app *signedTransferRequestService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *signedTransferRequestService) keynameByFromWalletID(fromWalletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_from_wallet_id:%s", app.keyname, fromWalletID.String())
}

func (app *signedTransferRequestService) retrievePartialSetByKeyname(keyname string, index int, amount int) (SignedTransferRequestPartialSet, error) {
	// retrieve transfer request uuids:
	ids := app.store.Sets().Retrieve(keyname, index, amount)

	// retrieve the transfer requests:
	reqs := []SignedTransferRequest{}
	for _, oneIDAsString := range ids {
		// cast the ID:
		id, idErr := uuid.FromString(oneIDAsString.(string))
		if idErr != nil {
			str := fmt.Sprintf("there is an element (%s) in the transfer request set (keyname: %s) that is not a valid TransferRequest UUID", keyname, oneIDAsString)
			return nil, errors.New(str)
		}

		req, reqErr := app.RetrieveByID(&id)
		if reqErr != nil {
			return nil, reqErr
		}

		reqs = append(reqs, req)
	}

	// retrieve the total amount of elements in the keyname:
	totalAmount := app.store.Sets().Len(keyname)

	// return:
	out := createSignedTransferPartialSet(reqs, index, totalAmount)
	return out, nil
}
