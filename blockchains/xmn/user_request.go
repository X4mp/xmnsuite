package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type storedUserRequest struct {
	ID       string `json:"id"`
	PubKey   string `json:"pubkey"`
	Shares   int    `json:"shares"`
	WalletID string `json:"wallet_id"`
}

func createStoredUserRequest(usr User) *storedUserRequest {
	out := storedUserRequest{
		ID:       usr.ID().String(),
		PubKey:   usr.PubKey().String(),
		Shares:   usr.Shares(),
		WalletID: usr.Wallet().ID().String(),
	}

	return &out
}

type userRequest struct {
	Usr User `json:"user"`
}

func createUserRequest(usr User) UserRequest {
	out := userRequest{
		Usr: usr,
	}

	return &out
}

// User returns the user
func (app *userRequest) User() User {
	return app.Usr
}

type userRequestPartialSet struct {
	Req   []UserRequest `json:"user_requests"`
	Idx   int           `json:"index"`
	TotAm int           `json:"total_amount"`
}

func createUserRequestPartialSet(req []UserRequest, index int, totalAmount int) UserRequestPartialSet {
	out := userRequestPartialSet{
		Req:   req,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out
}

// Users returns the users
func (app *userRequestPartialSet) Requests() []UserRequest {
	return app.Req
}

// Index returns the index
func (app *userRequestPartialSet) Index() int {
	return app.Idx
}

// Amount returns the amount
func (app *userRequestPartialSet) Amount() int {
	return len(app.Req)
}

// TotalAmount returns the total amount
func (app *userRequestPartialSet) TotalAmount() int {
	return app.TotAm
}

type userRequestService struct {
	keyname       string
	store         datastore.DataStore
	walletService WalletService
}

func createUserRequestService(store datastore.DataStore, walletService WalletService) UserRequestService {
	out := userRequestService{
		keyname:       "user_requests",
		store:         store,
		walletService: walletService,
	}

	return &out
}

// Save saves a UserRequest instance
func (app *userRequestService) Save(req UserRequest) error {
	// make sure the instance does not exists already:
	_, retErr := app.RetrieveByID(req.User().ID())
	if retErr == nil {
		return errors.New("the UserRequest instance already exists")
	}

	usr := req.User()
	usrID := usr.ID()
	walletID := usr.Wallet().ID()

	// make sure the wallet exists:
	_, retWalErr := app.walletService.RetrieveByID(walletID)
	if retWalErr != nil {
		str := fmt.Sprintf("there was an error while retrieving a Wallet instance by its ID (ID: %s): %s", walletID.String(), retWalErr.Error())
		return errors.New(str)
	}

	// make sure a UserRequest that contains both the WalletID and PubKey does not already exists:
	_, alreadyExistsReqErr := app.RetrieveByPubkeyAndWalletID(req.User().PubKey(), req.User().Wallet().ID())
	if alreadyExistsReqErr == nil {
		str := fmt.Sprintf("there is already a UserRequest that contains this PubKey (%s) and WalletID (%s)", req.User().PubKey().String(), req.User().Wallet().ID().String())
		return errors.New(str)
	}

	// create the set keys:
	keys := []string{
		app.keynameByWalletID(walletID),
		app.keynameByPubKey(req.User().PubKey()),
	}

	// add the userID to the set keynames:
	amountAddedToSets := app.store.Sets().AddMul(keys, usrID.String())

	if amountAddedToSets != 1 {
		// revert:
		app.store.Sets().DelMul(keys, usrID.String())

		// returns error:
		str := fmt.Sprintf("there was an error while adding the UserRequest (UserID: %s) to the sets... reverting", usrID.String())
		return errors.New(str)
	}

	// save the object:
	keyname := app.keynameByID(usrID)
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: keyname,
		Obj: createStoredUserRequest(usr),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the UserRequest instance")
	}

	return nil
}

// Delete deletes a user request
func (app *userRequestService) Delete(req UserRequest) error {
	userID := req.User().ID()
	keynameByID := app.keynameByID(userID)
	amountDelObjects := app.store.Objects().Keys().Delete(keynameByID)
	if amountDelObjects != 1 {
		str := fmt.Sprintf("there was an error while deleting the UserRequest (UserID: %s), key: %s", req.User().ID().String(), keynameByID)
		return errors.New(str)
	}

	// create the set keys:
	keys := []string{
		app.keynameByWalletID(req.User().Wallet().ID()),
		app.keynameByPubKey(req.User().PubKey()),
	}

	// delete the userID from the set keynames:
	amountDeletedFromSets := app.store.Sets().DelMul(keys, userID.String())

	if amountDeletedFromSets != 1 {
		// revert:
		app.store.Sets().DelMul(keys, userID.String())

		// returns error:
		str := fmt.Sprintf("there was an error while deleting the UserRequest (UserID: %s) from the sets... reverting", userID.String())
		return errors.New(str)
	}

	return nil
}

// RetrieveByID retrieves a UserRequest instance by its ID
func (app *userRequestService) RetrieveByID(id *uuid.UUID) (UserRequest, error) {
	keyname := app.keynameByID(id)
	obj := objects.ObjInKey{
		Key: keyname,
		Obj: new(storedUserRequest),
	}

	amount := app.store.Objects().Retrieve(&obj)
	if amount != 1 {
		str := fmt.Sprintf("there was an error while retrieving the UserRequest (ID: %s)", id.String())
		return nil, errors.New(str)
	}

	if req, ok := obj.Obj.(*storedUserRequest); ok {
		return app.FromStoredToUserRequest(req)
	}

	return nil, errors.New("the retrieved data cannot be casted to a UserRequest instance")
}

// RetrieveByPubkeyAndWalletID retrieves a UserRequest instance by uts PubKey and WalletID
func (app *userRequestService) RetrieveByPubkeyAndWalletID(pubKey crypto.PublicKey, walletID *uuid.UUID) (UserRequest, error) {
	// retireve keynames:
	keynameByWalletID := app.keynameByWalletID(walletID)
	keynameByPubKey := app.keynameByPubKey(pubKey)

	// intersect:
	ids := app.store.Sets().Inter(keynameByWalletID, keynameByPubKey)
	if len(ids) <= 0 {
		str := fmt.Sprintf("there is no UserRequest that contains both that PubKey (%s) and that WalletID (%s)", pubKey.String(), walletID.String())
		return nil, errors.New(str)
	}

	if len(ids) == 1 {
		// cast the ID:
		id, idErr := uuid.FromString(ids[0].(string))
		if idErr != nil {
			str := fmt.Sprintf("the element stored in the set is not a valid UUID: %s", idErr.Error())
			return nil, errors.New(str)
		}

		// retrieve the instance, then return it:
		return app.RetrieveByID(&id)
	}

	str := fmt.Sprintf("there is %d UserRequest instances that contains both that PubKey (%s) and that WalletID (%s), this should never happen", len(ids), pubKey.String(), walletID.String())
	return nil, errors.New(str)
}

// RetrieveByWalletID retrieve []UserRequest by its walletID
func (app *userRequestService) RetrieveByWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestPartialSet, error) {
	keyname := app.keynameByWalletID(walletID)
	return app.retrieveByKeyname(keyname, index, amount)
}

// RetrieveByPubKey retrieve []UserRequest by its pubKey
func (app *userRequestService) RetrieveByPubKey(pubKey crypto.PublicKey, index int, amount int) (UserRequestPartialSet, error) {
	keyname := app.keynameByPubKey(pubKey)
	return app.retrieveByKeyname(keyname, index, amount)
}

// FromStoredToUserRequest converts a stored UserRequest to a UserRequest instance
func (app *userRequestService) FromStoredToUserRequest(req *storedUserRequest) (UserRequest, error) {
	// cast the requestID:
	reqID, reqIDErr := uuid.FromString(req.ID)
	if reqIDErr != nil {
		str := fmt.Sprintf("the ID (%s) is invalid", req.ID)
		return nil, errors.New(str)
	}

	// cast the pubKey:
	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: req.PubKey,
	})

	// cast the walletID:
	walletID, walletIDErr := uuid.FromString(req.WalletID)
	if walletIDErr != nil {
		str := fmt.Sprintf("the walletID (%s) is invalid", req.WalletID)
		return nil, errors.New(str)
	}

	// retrieve the wallet:
	wal, walErr := app.walletService.RetrieveByID(&walletID)
	if walErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the Wallet instance (ID: %s): %s", req.WalletID, walErr.Error())
		return nil, errors.New(str)
	}

	// create user:
	usr := createUser(&reqID, pubKey, req.Shares, wal)

	// create request:
	usrReq := createUserRequest(usr)
	return usrReq, nil
}

func (app *userRequestService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *userRequestService) keynameByWalletID(walletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_wallet_id:%s", app.keyname, walletID.String())
}

func (app *userRequestService) keynameByPubKey(pubKey crypto.PublicKey) string {
	return fmt.Sprintf("%s:by_pubkey:%s", app.keyname, pubKey.String())
}

func (app *userRequestService) retrieveByKeyname(keyname string, index int, amount int) (UserRequestPartialSet, error) {
	reqs := []UserRequest{}
	uncastedUserIDs := app.store.Sets().Retrieve(keyname, index, amount)
	for _, oneUncastedUserID := range uncastedUserIDs {
		userID, userIDErr := uuid.FromString(oneUncastedUserID.(string))
		if userIDErr != nil {
			str := fmt.Sprintf("one of the elements in the set (key: %s) is not a valid userID (element: %s): %s", keyname, oneUncastedUserID.(string), userIDErr.Error())
			return nil, errors.New(str)
		}

		oneReq, oneReqErr := app.RetrieveByID(&userID)
		if oneReqErr != nil {
			return nil, oneReqErr
		}

		reqs = append(reqs, oneReq)
	}

	// retireve the total amount:
	totAmount := app.store.Sets().Len(keyname)
	ps := createUserRequestPartialSet(reqs, index, totAmount)
	return ps, nil
}
