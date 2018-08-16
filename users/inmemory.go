package users

import (
	"errors"
	"fmt"

	"github.com/XMNBlockchain/datamint/objects"
	crypto "github.com/tendermint/tendermint/crypto"
)

type concreteUsers struct {
	store objects.Objects
}

func createConcreteUsers() Users {
	out := concreteUsers{
		store: objects.SDKFunc.Create(),
	}

	return &out
}

// Head returns the head
func (app *concreteUsers) Objects() objects.Objects {
	return app.store
}

// Exists returns true if the user exists, false otherwise
func (app *concreteUsers) Exists(pubKey crypto.PubKey) bool {
	key := app.genKey(pubKey)
	return app.store.Keys().Exists(key) == 1
}

// Add adds a user
func (app *concreteUsers) Insert(pubKey crypto.PubKey) error {
	if app.Exists(pubKey) {
		str := fmt.Sprintf("the given public key (%s) is already assigned to another user", pubKey)
		return errors.New(str)
	}

	key := app.genKey(pubKey)
	app.store.Save(&objects.ObjInKey{
		Key: key,
		Obj: pubKey,
	})
	return nil
}

// Delete deletes a user
func (app *concreteUsers) Delete(pubKey crypto.PubKey) error {
	if !app.Exists(pubKey) {
		str := fmt.Sprintf("the given public key (%s) is not assigned to a user", pubKey)
		return errors.New(str)
	}

	key := app.genKey(pubKey)
	app.store.Keys().Delete(key)
	return nil
}

func (app *concreteUsers) genKey(pubKey crypto.PubKey) string {
	return fmt.Sprintf("user:by_pubkey:%s", pubKey)
}
