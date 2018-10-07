package datastore

import (
	"github.com/xmnservices/xmnsuite/hashtree"
	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/lists"
	"github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/roles"
	"github.com/xmnservices/xmnsuite/users"
)

/*
 * DataStore
 *
 */

type concreteDataStore struct {
	K    keys.Keys
	L    lists.Lists
	S    lists.Lists
	Objs objects.Objects
	Usrs users.Users
	Rols roles.Roles
}

func createConcreteDataStore() DataStore {
	out := concreteDataStore{
		K:    keys.SDKFunc.Create(),
		L:    lists.SDKFunc.CreateList(),
		S:    lists.SDKFunc.CreateSet(),
		Objs: objects.SDKFunc.Create(),
		Usrs: users.SDKFunc.Create(),
		Rols: roles.SDKFunc.Create(),
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
			app.Objects().Keys().Head().Head().Get(),
			app.Users().Objects().Keys().Head().Head().Get(),
			app.Roles().Lists().Objects().Keys().Head().Head().Get(),
		},
	})

	return head
}

// Copy copies the datastore
func (app *concreteDataStore) Copy() DataStore {
	ck := app.K.Copy()
	cl := app.L.Copy()
	cs := app.S.Copy()
	cobjs := app.Objs.Copy()
	usrs := app.Usrs.Copy()
	rols := app.Rols.Copy()
	out := concreteDataStore{
		K:    ck,
		L:    cl,
		S:    cs,
		Objs: cobjs,
		Usrs: usrs,
		Rols: rols,
	}

	return &out
}

// Keys returns the keys datastore
func (app *concreteDataStore) Keys() keys.Keys {
	return app.K
}

// Lists returns the lists datastore
func (app *concreteDataStore) Lists() lists.Lists {
	return app.L
}

// Sets returns the sets datastore
func (app *concreteDataStore) Sets() lists.Lists {
	return app.S
}

// Objects returns the objects datastore
func (app *concreteDataStore) Objects() objects.Objects {
	return app.Objs
}

// Users returns the users datastore
func (app *concreteDataStore) Users() users.Users {
	return app.Usrs
}

// Roles returns the roles datastore
func (app *concreteDataStore) Roles() roles.Roles {
	return app.Rols
}
