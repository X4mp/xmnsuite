package users

import (
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

// Objects returns the objects
func (app *concreteUsers) Objects() objects.Objects {
	return app.store
}

// Copy copues the Users instance
func (app *concreteUsers) Copy() Users {
	out := concreteUsers{
		store: app.store.Copy(),
	}

	return &out
}

// Key returns the key where the user is stored
func (app *concreteUsers) Key(pubKey crypto.PubKey) string {
	return fmt.Sprintf("user:by_pubkey:%s", pubKey)
}

// Exists returns true if the user exists, false otherwise
func (app *concreteUsers) Exists(pubKey crypto.PubKey) bool {
	key := app.Key(pubKey)
	return app.store.Keys().Exists(key) == 1
}

// Add adds a user
func (app *concreteUsers) Insert(pubKey crypto.PubKey) bool {
	if app.Exists(pubKey) {
		return false
	}

	key := app.Key(pubKey)
	app.store.Save(&objects.ObjInKey{
		Key: key,
		Obj: pubKey,
	})

	return true
}

// Delete deletes a user
func (app *concreteUsers) Delete(pubKey crypto.PubKey) bool {
	if !app.Exists(pubKey) {
		return false
	}

	key := app.Key(pubKey)
	app.store.Keys().Delete(key)
	return true
}
