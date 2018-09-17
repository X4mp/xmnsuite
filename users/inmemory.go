package users

import (
	"fmt"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/objects"
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
func (app *concreteUsers) Key(pubKey crypto.PublicKey) string {
	return fmt.Sprintf("user:by_pubkey:%s", pubKey)
}

// Exists returns true if the user exists, false otherwise
func (app *concreteUsers) Exists(pubKey crypto.PublicKey) bool {
	key := app.Key(pubKey)
	return app.store.Keys().Exists(key) == 1
}

// Add adds a user
func (app *concreteUsers) Insert(pubKey crypto.PublicKey) bool {
	if app.Exists(pubKey) {
		return false
	}

	key := app.Key(pubKey)
	app.store.Save(&objects.ObjInKey{
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
	app.store.Keys().Delete(key)
	return true
}
