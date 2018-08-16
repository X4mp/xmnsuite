package users

import (
	"errors"
	"fmt"

	"github.com/XMNBlockchain/datamint/hashtree"
	"github.com/XMNBlockchain/datamint/keys"
	crypto "github.com/tendermint/tendermint/crypto"
)

type concreteUsers struct {
	store keys.Keys
}

func createConcreteUsers() Users {
	out := concreteUsers{
		store: keys.SDKFunc.Create(),
	}

	return &out
}

// Head returns the head
func (app *concreteUsers) Head() hashtree.HashTree {
	return app.store.Head()
}

// Len returns the amount of users
func (app *concreteUsers) Len() int {
	return app.store.Len()
}

// Exists returns true if the user exists, false otherwise
func (app *concreteUsers) Exists(pubKey crypto.PubKey) bool {
	key := app.genKey(pubKey)
	return app.store.Exists(key) == 1
}

// Add adds a user
func (app *concreteUsers) Insert(pubKey crypto.PubKey) error {
	key := app.genKey(pubKey)
	if app.store.Exists(key) == 1 {
		str := fmt.Sprintf("the given public key (%s) is already assigned to another user", key)
		return errors.New(str)
	}

	app.store.Save(key, pubKey)
	return nil
}

// Delete deletes a user
func (app *concreteUsers) Delete(pubKey crypto.PubKey) error {
	key := app.genKey(pubKey)
	if app.store.Exists(key) != 1 {
		str := fmt.Sprintf("the given public key (%s) is not assigned to a user", key)
		return errors.New(str)
	}

	app.store.Delete(key)
	return nil
}

func (app *concreteUsers) genKey(pubKey crypto.PubKey) string {
	return fmt.Sprintf("user:by_pubkey:%s", pubKey)
}
