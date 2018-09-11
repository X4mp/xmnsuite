package datastore

import (
	"github.com/XMNBlockchain/xmnsuite/hashtree"
	"github.com/XMNBlockchain/xmnsuite/keys"
	"github.com/XMNBlockchain/xmnsuite/lists"
	"github.com/XMNBlockchain/xmnsuite/objects"
	"github.com/XMNBlockchain/xmnsuite/roles"
	"github.com/XMNBlockchain/xmnsuite/users"
)

// DataStore represents the datastore
type DataStore interface {
	Head() hashtree.HashTree
	Copy() DataStore
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
