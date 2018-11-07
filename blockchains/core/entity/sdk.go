package entity

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

// ToStorable represents the ToStorable func type
type ToStorable func(ins Entity) (interface{}, error)

// ToEntity represents the ToEntity func type
type ToEntity func(rep Repository, data interface{}) (Entity, error)

// Keynames returns the keynames related to the entity
type Keynames func(ins Entity) ([]string, error)

// Sync syncs the sub entities with the database.  Can be used to store the sub entities in the database, before storing the current entity
type Sync func(rep Repository, service Service, ins Entity) error

// Entity represents an entity instance
type Entity interface {
	ID() *uuid.UUID
}

// MetaData represents the metadata
type MetaData interface {
	Name() string
	Keyname() string
	ToEntity() ToEntity
	CopyStorable() interface{}
}

// Representation represents an entity representation
type Representation interface {
	MetaData() MetaData
	ToStorable() ToStorable
	HasKeynames() bool
	Keynames() Keynames
	HasSync() bool
	Sync() Sync
}

// PartialSet represents an  entity partial set
type PartialSet interface {
	Instances() []Entity
	Index() int
	Amount() int
	TotalAmount() int
}

// Service represents an entity service
type Service interface {
	Save(ins Entity, rep Representation) error
	Delete(ins Entity, rep Representation) error
}

// Repository represents an entity repository
type Repository interface {
	RetrieveByID(met MetaData, id *uuid.UUID) (Entity, error)
	RetrieveByIntersectKeynames(met MetaData, keynames []string) (Entity, error)
	RetrieveSetByKeyname(met MetaData, keyname string, index int, amount int) (PartialSet, error)
	RetrieveSetByIntersectKeynames(met MetaData, keynames []string, index int, amount int) (PartialSet, error)
}

// Controllers represents the func controllers
type Controllers interface {
	Save() routers.SaveTransactionFn
	Delete() routers.DeleteTransactionFn
	RetrieveByID() routers.QueryFn
	RetrieveByIntersectKeynames() routers.QueryFn
	RetrieveSetByIntersectKeynames() routers.QueryFn
	RetrieveSetByKeyname() routers.QueryFn
}

// CreateMetaDataParams represents the MetaData params
type CreateMetaDataParams struct {
	Name          string
	ToEntity      ToEntity
	EmptyStorable interface{}
}

// CreateRepresentationParams represents the Representation params
type CreateRepresentationParams struct {
	Met        MetaData
	ToStorable ToStorable
	Keynames   Keynames
	Sync       Sync
}

// CreateRepositoryParams represents the Repository params
type CreateRepositoryParams struct {
	Store datastore.DataStore
}

// CreateServiceParams represents the Service params
type CreateServiceParams struct {
	Store datastore.DataStore
}

// CreateSDKRepositoryParams represents the CreateSDKRepository params
type CreateSDKRepositoryParams struct {
	PK     crypto.PrivateKey
	Client applications.Client
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK     crypto.PrivateKey
	Client applications.Client
}

// CreateControllersParams represents the Controllers params
type CreateControllersParams struct {
	Met                      MetaData
	Rep                      Representation
	DefaultAmountOfElements  int
	GazPricePerKb            int
	OverwriteIfAlreadyExists bool
	RouterRoleKey            string
}

// SDKFunc represents the Entity SDK func
var SDKFunc = struct {
	CreateMetaData       func(params CreateMetaDataParams) MetaData
	CreateRepresentation func(params CreateRepresentationParams) Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateService        func(params CreateServiceParams) Service
	CreateControllers    func(params CreateControllersParams) Controllers
	CreateSDKRepository  func(params CreateSDKRepositoryParams) Repository
	CreateSDKService     func(params CreateSDKServiceParams) Service
}{
	CreateMetaData: func(params CreateMetaDataParams) MetaData {
		met, metErr := createMetaData(params.Name, params.ToEntity, params.EmptyStorable)
		if metErr != nil {
			panic(metErr)
		}

		return met
	},
	CreateRepresentation: func(params CreateRepresentationParams) Representation {
		if params.Keynames != nil && params.Sync != nil {
			out := createRepresentationWithKeynamesAndSync(params.Met, params.ToStorable, params.Keynames, params.Sync)
			return out
		}

		if params.Keynames != nil {
			out := createRepresentationWithKeynames(params.Met, params.ToStorable, params.Keynames)
			return out
		}

		if params.Sync != nil {
			out := createRepresentationWithSync(params.Met, params.ToStorable, params.Sync)
			return out
		}

		out := createRepresentation(params.Met, params.ToStorable)
		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		out := createRepository(params.Store)
		return out
	},
	CreateService: func(params CreateServiceParams) Service {
		rep := createRepository(params.Store)
		out := createService(params.Store, rep)
		return out
	},
	CreateControllers: func(params CreateControllersParams) Controllers {
		out := createControllers(
			params.Met,
			params.Rep,
			params.DefaultAmountOfElements,
			params.GazPricePerKb,
			params.OverwriteIfAlreadyExists,
			params.RouterRoleKey,
		)

		return out
	},
	CreateSDKRepository: func(params CreateSDKRepositoryParams) Repository {
		out := createSDKRepository(params.PK, params.Client)
		return out
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client)
		return out
	},
}
