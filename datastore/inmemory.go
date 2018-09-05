package datastore

import (
	"github.com/XMNBlockchain/datamint/hashtree"
	"github.com/XMNBlockchain/datamint/keys"
	"github.com/XMNBlockchain/datamint/lists"
	"github.com/XMNBlockchain/datamint/objects"
	"github.com/XMNBlockchain/datamint/roles"
	"github.com/XMNBlockchain/datamint/users"
)

/*
 * DataStore
 *
 */

type concreteDataStore struct {
	k    keys.Keys
	l    lists.Lists
	s    lists.Lists
	objs objects.Objects
	usrs users.Users
	rols roles.Roles
}

func createConcreteDataStore() DataStore {
	out := concreteDataStore{
		k:    keys.SDKFunc.Create(),
		l:    lists.SDKFunc.CreateList(),
		s:    lists.SDKFunc.CreateSet(),
		objs: objects.SDKFunc.Create(),
		usrs: users.SDKFunc.Create(),
		rols: roles.SDKFunc.Create(),
	}

	return &out
}

// Head returns the hashtree of the datastore
func (app *concreteDataStore) Head() hashtree.HashTree {
	head := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			app.Keys().Head().Head().Get(),
			app.Lists().Objects().Keys().Head().Head().Get(),
			app.Sets().Objects().Keys().Head().Head().Get(),
			app.objs.Keys().Head().Head().Get(),
			app.usrs.Objects().Keys().Head().Head().Get(),
			app.rols.Lists().Objects().Keys().Head().Head().Get(),
		},
	})

	return head
}

// Copy copies the datastore
func (app *concreteDataStore) Copy() DataStore {
	ck := app.k.Copy()
	cobjs := app.objs.Copy()
	out := concreteDataStore{
		k:    ck,
		l:    app.Lists(),
		s:    app.Sets(),
		objs: cobjs,
		usrs: app.Users(),
		rols: app.Roles(),
	}

	return &out
}

// Keys returns the keys datastore
func (app *concreteDataStore) Keys() keys.Keys {
	return app.k
}

// Lists returns the lists datastore
func (app *concreteDataStore) Lists() lists.Lists {
	return app.l
}

// Sets returns the sets datastore
func (app *concreteDataStore) Sets() lists.Lists {
	return app.s
}

// Objects returns the objects datastore
func (app *concreteDataStore) Objects() objects.Objects {
	return app.objs
}

// Users returns the users datastore
func (app *concreteDataStore) Users() users.Users {
	return app.usrs
}

// Roles returns the roles datastore
func (app *concreteDataStore) Roles() roles.Roles {
	return app.rols
}
