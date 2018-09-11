package datastore

import (
	"github.com/xmnservices/xmnsuite/hashtree"
	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/lists"
	"github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/roles"
	"github.com/xmnservices/xmnsuite/users"
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
