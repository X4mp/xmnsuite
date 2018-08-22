package datastore

import (
	"github.com/XMNBlockchain/datamint/hashtree"
	"github.com/XMNBlockchain/datamint/keys"
	"github.com/XMNBlockchain/datamint/lists"
	"github.com/XMNBlockchain/datamint/objects"
	"github.com/XMNBlockchain/datamint/roles"
	"github.com/XMNBlockchain/datamint/users"
)

// DataStore represents the datastore
type DataStore interface {
	Head() hashtree.HashTree
	Keys() keys.Keys
	Lists() lists.Lists
	Sets() lists.Lists
	Objects() objects.Objects
	Users() users.Users
	Roles() roles.Roles
}

// SDKFunc represents the datastore SDK func
var SDKFunc = struct {
	Create func() DataStore
}{
	Create: func() DataStore {
		return createConcreteDataStore()
	},
}
