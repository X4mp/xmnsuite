package datastore

import (
	"path/filepath"

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

// StoredDataStore represents a stored datastore
type StoredDataStore interface {
	DataStore() DataStore
	Save() error
}

// ServiceParams represents the service params
type ServiceParams struct {
	DirPath string
}

// StoredDataStoreParams represents the StoredDataStore params
type StoredDataStoreParams struct {
	FilePath string
}

// SDKFunc represents the datastore SDK func
var SDKFunc = struct {
	Create                func() DataStore
	CreateService         func(params ServiceParams) Service
	CreateStoredDataStore func(params StoredDataStoreParams) StoredDataStore
}{
	Create: func() DataStore {
		return createConcreteDataStore()
	},
	CreateService: func(params ServiceParams) Service {
		return createFileService(params.DirPath)
	},
	CreateStoredDataStore: func(params StoredDataStoreParams) StoredDataStore {
		dirPath := filepath.Dir(params.FilePath)
		fileName := filepath.Base(params.FilePath)
		serv := createFileService(dirPath)
		ds, dsErr := serv.Retrieve(fileName)
		if dsErr != nil {
			ds = createConcreteDataStore()
		}

		st := createConcreteStoredDataStore(ds, serv, fileName)
		return st
	},
}
