package users

import (
	"fmt"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type concreteUsers struct {
	Store objects.Objects
}

func createConcreteUsers() Users {
	out := concreteUsers{
		Store: objects.SDKFunc.Create(),
	}

	return &out
}

// Objects returns the objects
func (app *concreteUsers) Objects() objects.Objects {
	return app.Store
}

// Copy copues the Users instance
func (app *concreteUsers) Copy() Users {
	out := concreteUsers{
		Store: app.Store.Copy(),
	}

	return &out
}

// Key returns the key where the user is stored
func (app *concreteUsers) Key(pubKey crypto.PublicKey) string {
	return fmt.Sprintf("user:by_pubkey:%s", pubKey)
}

// Exists returns true if the user exists, false otherwise
func (app *concreteUsers) Exists(pubKey crypto.PublicKey) bool {
	key := app.Key(pubKey)
	return app.Store.Keys().Exists(key) == 1
}

// Add adds a user
func (app *concreteUsers) Insert(pubKey crypto.PublicKey) bool {
	if app.Exists(pubKey) {
		return false
	}

	key := app.Key(pubKey)
	app.Store.Save(&objects.ObjInKey{
		Key: key,
		Obj: pubKey.String(),
	})

	return true
}

// Delete deletes a user
func (app *concreteUsers) Delete(pubKey crypto.PublicKey) bool {
	if !app.Exists(pubKey) {
		return false
	}

	key := app.Key(pubKey)
	app.Store.Keys().Delete(key)
	return true
}
