package datastore

import (
	"github.com/xmnservices/xmnsuite/datastore/keys"
	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/datastore/objects"
	"github.com/xmnservices/xmnsuite/datastore/roles"
	"github.com/xmnservices/xmnsuite/datastore/users"
	"github.com/xmnservices/xmnsuite/hashtree"
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

// Service saves and retrieves datastores
type Service interface {
	Save(ds DataStore, filePath string) error
	Retrieve(filePath string) (DataStore, error)
}

// SDKFunc represents the datastore SDK func
var SDKFunc = struct {
	Create func() DataStore
}{
	Create: func() DataStore {
		return createConcreteDataStore()
	},
}
