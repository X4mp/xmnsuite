package xmn

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

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

type storedInitialDeposit struct {
	ToWalletID string `json:"to_wallet_id"`
	Am         int    `json:"amount"`
}

func createStoredInitialDeposit(dep InitialDeposit) *storedInitialDeposit {
	out := storedInitialDeposit{
		ToWalletID: dep.To().ID().String(),
		Am:         dep.Amount(),
	}

	return &out
}

type initialDepositService struct {
	keyname       string
	store         datastore.DataStore
	walletService WalletService
}

func createInitialDepositService(store datastore.DataStore, walletService WalletService) InitialDepositService {
	out := initialDepositService{
		keyname:       "initial-deposit",
		store:         store,
		walletService: walletService,
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
	saveWalErr := app.walletService.Save(initialDep.To())
	if saveWalErr != nil {
		str := fmt.Sprintf("there was an error while saving the Wallet instance, in the InitialDeposit instance: %s", saveWalErr.Error())
		return errors.New(str)
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keyname,
		Obj: storedInitialDeposit{
			ToWalletID: initialDep.To().ID().String(),
			Am:         initialDep.Amount(),
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
		walID, walIDErr := uuid.FromString(storedTok.ToWalletID)
		if walIDErr != nil {
			return nil, walIDErr
		}

		// retrieve the wallet:
		wal, walErr := app.walletService.RetrieveByID(&walID)
		if walErr != nil {
			return nil, walErr
		}

		out := createInitialDeposit(wal, storedTok.Am)
		return out, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a InitialDeposit instance")
}
