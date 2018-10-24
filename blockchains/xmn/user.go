package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type user struct {
	UUID  *uuid.UUID       `json:"id"`
	PKey  crypto.PublicKey `json:"pubkey"`
	Shres int              `json:"shares"`
	Wal   Wallet           `json:"wallet"`
}

type storedUser struct {
	ID       string `json:"id"`
	PubKey   string `json:"pubkey"`
	Shares   int    `json:"shares"`
	WalletID string `json:"wallet_id"`
}

func createUser(id *uuid.UUID, pubKey crypto.PublicKey, shares int, wallet Wallet) User {
	out := user{
		UUID:  id,
		PKey:  pubKey,
		Shres: shares,
		Wal:   wallet,
	}

	return &out
}

// ID returns the ID
func (app *user) ID() *uuid.UUID {
	return app.UUID
}

// PubKey returns the PublicKey
func (app *user) PubKey() crypto.PublicKey {
	return app.PKey
}

// Shares returns the shares
func (app *user) Shares() int {
	return app.Shres
}

// Wallet returns the wallet
func (app *user) Wallet() Wallet {
	return app.Wal
}

type userPartialSet struct {
	Usrs  []User `json:"users"`
	Idx   int    `json:"index"`
	TotAm int    `json:"total_amount"`
}

func createUserPartialSet(usrs []User, index int, totalAmount int) UserPartialSet {
	out := userPartialSet{
		Usrs:  usrs,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out
}

// Users returns the users
func (app *userPartialSet) Users() []User {
	return app.Usrs
}

// Index returns the index
func (app *userPartialSet) Index() int {
	return app.Idx
}

// Amount returns the amount
func (app *userPartialSet) Amount() int {
	return len(app.Usrs)
}

// TotalAmount returns the total amount
func (app *userPartialSet) TotalAmount() int {
	return app.TotAm
}

type userService struct {
	keyname       string
	store         datastore.DataStore
	walletService WalletService
}

func createUserService(store datastore.DataStore, walletService WalletService) UserService {
	out := userService{
		keyname:       "users",
		store:         store,
		walletService: walletService,
	}

	return &out
}

// Save saves a User instance
func (app *userService) Save(usr User) error {
	// make sure the user does not already exists:
	_, retUsrErr := app.RetrieveByID(usr.ID())
	if retUsrErr == nil {
		str := fmt.Sprintf("the User (ID: %s) already exists", usr.ID().String())
		return errors.New(str)
	}

	// create the set keys:
	keys := []string{
		app.keynameByWalletID(usr.Wallet().ID()),
		app.keynameByPubKey(usr.PubKey()),
	}

	// add the ID to the set keynames:
	amountAddedToSets := app.store.Sets().AddMul(keys, usr.ID().String())

	if amountAddedToSets != 1 {
		// revert:
		app.store.Sets().DelMul(keys, usr.ID().String())

		// returns error:
		str := fmt.Sprintf("there was an error while adding the User (ID: %s) to the sets... reverting", usr.ID().String())
		return errors.New(str)
	}

	// save the object:
	keyname := app.keynameByID(usr.ID())
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: keyname,
		Obj: storedUser{
			ID:       usr.ID().String(),
			PubKey:   usr.PubKey().String(),
			Shares:   usr.Shares(),
			WalletID: usr.Wallet().ID().String(),
		},
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the User instance")
	}

	return nil
}

// RetrieveByID retrieves a user by its ID
func (app *userService) RetrieveByID(id *uuid.UUID) (User, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keynameByID(id),
		Obj: new(storedUser),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		str := fmt.Sprintf("there was an error while retrieving the User instance (ID: %s)", id.String())
		return nil, errors.New(str)
	}

	// cast the instance:
	if storedUser, ok := obj.Obj.(*storedUser); ok {
		return app.fromStoredUserToUser(storedUser)
	}

	return nil, errors.New("the retrieved data cannot be casted to a User instance")
}

// RetrieveByWalletID retrieves users from a walletID
func (app *userService) RetrieveByWalletID(walletID *uuid.UUID, index int, amount int) (UserPartialSet, error) {
	keyname := app.keynameByWalletID(walletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

// RetrieveRetrieveByPubKeyByID retrieves users from pubKey
func (app *userService) RetrieveByPubKey(pubKey crypto.PublicKey, index int, amount int) (UserPartialSet, error) {
	keyname := app.keynameByPubKey(pubKey)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

func (app *userService) fromStoredUserToUser(storedUsr *storedUser) (User, error) {
	// cast the walletID:
	walletID, walletIDErr := uuid.FromString(storedUsr.WalletID)
	if walletIDErr != nil {
		return nil, walletIDErr
	}

	// retrieve the wallet:
	retWallet, retWalletErr := app.walletService.RetrieveByID(&walletID)
	if retWalletErr != nil {
		return nil, retWalletErr
	}

	// cast the ID:
	id, idErr := uuid.FromString(storedUsr.ID)
	if idErr != nil {
		return nil, idErr
	}

	//cast the pubkey:
	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: storedUsr.PubKey,
	})

	// create user:
	usr := createUser(&id, pubKey, storedUsr.Shares, retWallet)
	return usr, nil
}

func (app *userService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *userService) keynameByWalletID(walletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_wallet_id:%s", app.keyname, walletID.String())
}

func (app *userService) keynameByPubKey(pubKey crypto.PublicKey) string {
	return fmt.Sprintf("%s:by_pubkey:%s", app.keyname, pubKey.String())
}

func (app *userService) retrievePartialSetByKeyname(keyname string, index int, amount int) (UserPartialSet, error) {
	// retrieve user uuids:
	ids := app.store.Sets().Retrieve(keyname, index, amount)

	// retrieve the users:
	users := []User{}
	for _, oneIDAsString := range ids {
		// cast the ID:
		id, idErr := uuid.FromString(oneIDAsString.(string))
		if idErr != nil {
			str := fmt.Sprintf("there is an element (%s) in the users set (keyname: %s) that is not a valid users UUID", keyname, oneIDAsString)
			return nil, errors.New(str)
		}

		usr, usrErr := app.RetrieveByID(&id)
		if usrErr != nil {
			return nil, usrErr
		}

		users = append(users, usr)
	}

	// retrieve the total amount of elements in the keyname:
	totalAmount := app.store.Sets().Len(keyname)

	// return:
	out := createUserPartialSet(users, index, totalAmount)
	return out, nil
}
