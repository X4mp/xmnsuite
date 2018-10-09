package datastore

import (
	"github.com/xmnservices/xmnsuite/datastore/keys"
	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/datastore/objects"
	"github.com/xmnservices/xmnsuite/datastore/roles"
	"github.com/xmnservices/xmnsuite/datastore/users"
	"github.com/xmnservices/xmnsuite/hashtree"
)

/*
 * StoredDataStore
 *
 */

type concreteStoredDataStore struct {
	ds       DataStore
	serv     Service
	fileName string
}

func createConcreteStoredDataStore(ds DataStore, serv Service, fileName string) StoredDataStore {
	out := concreteStoredDataStore{
		ds:       ds,
		serv:     serv,
		fileName: fileName,
	}

	return &out
}

// DataStore returns the datastore
func (app *concreteStoredDataStore) DataStore() DataStore {
	return app.ds
}

// Save save the DataStore on disk
func (app *concreteStoredDataStore) Save() error {
	saveErr := app.serv.Save(app.ds, app.fileName)
	if saveErr != nil {
		return saveErr
	}

	return nil
}

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
