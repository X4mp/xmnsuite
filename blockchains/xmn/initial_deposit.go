package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type initialDeposit struct {
	ToUser User `json:"to"`
	Am     int  `json:"amount"`
}

func createInitialDeposit(toUsr User, amount int) InitialDeposit {
	out := initialDeposit{
		ToUser: toUsr,
		Am:     amount,
	}

	return &out
}

// To returns the to user
func (app *initialDeposit) To() User {
	return app.ToUser
}

// Amount returns the amount
func (app *initialDeposit) Amount() int {
	return app.Am
}

type storedInitialDeposit struct {
	ToUserID string `json:"to_user_id"`
	Am       int    `json:"amount"`
}

type initialDepositService struct {
	keyname       string
	store         datastore.DataStore
	walletService WalletService
	userService   UserService
}

func createInitialDepositService(store datastore.DataStore, walletService WalletService, userService UserService) InitialDepositService {
	out := initialDepositService{
		keyname:       "initial-deposit",
		store:         store,
		walletService: walletService,
		userService:   userService,
	}

	return &out
}

// Save saves the InitialDeposit instance
func (app *initialDepositService) Save(initialDep InitialDeposit) error {
	// make sure the instance does not exists already:
	_, retErr := app.Retrieve()
	if retErr == nil {
		return errors.New("the InitialDeposit instance already exists")
	}

	// save the wallet:
	saveWalletErr := app.walletService.Save(initialDep.To().Wallet())
	if saveWalletErr != nil {
		str := fmt.Sprintf("there was an error while saving the Wallet instance, in the User instance, in the InitialDeposit instance: %s", saveWalletErr.Error())
		return errors.New(str)
	}

	// save the user:
	saveUserErr := app.userService.Save(initialDep.To())
	if saveUserErr != nil {
		str := fmt.Sprintf("there was an error while saving the User instance, in the InitialDeposit instance: %s", saveUserErr.Error())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keyname,
		Obj: storedInitialDeposit{
			ToUserID: initialDep.To().ID().String(),
			Am:       initialDep.Amount(),
		},
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the InitialDeposit instance")
	}

	return nil
}

// Retrieve retrieves the InitialDeposit instance
func (app *initialDepositService) Retrieve() (InitialDeposit, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keyname,
		Obj: new(storedInitialDeposit),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		return nil, errors.New("there was an error while retrieving the InitialDeposit instance")
	}

	// cast the instance:
	if storedTok, ok := obj.Obj.(*storedInitialDeposit); ok {
		// cast the ID:
		userID, userIDErr := uuid.FromString(storedTok.ToUserID)
		if userIDErr != nil {
			return nil, userIDErr
		}

		// retrieve the wallet:
		usr, usrErr := app.userService.RetrieveByID(&userID)
		if usrErr != nil {
			return nil, usrErr
		}

		out := createInitialDeposit(usr, storedTok.Am)
		return out, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a InitialDeposit instance")
}
