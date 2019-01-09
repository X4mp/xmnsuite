package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Request represents an entity request
type Request interface {
	ID() *uuid.UUID
	From() user.User
	HasSave() bool
	Save() entity.Entity
	HasDelete() bool
	Delete() entity.Entity
	Reason() string
	Keyname() keyname.Keyname
}

// Service represents an entity service
type Service interface {
	Save(req Request, entityRep entity.Representation) error
}

// Normalized represents a normalized request
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID           *uuid.UUID
	FromUser     user.User
	SaveEntity   entity.Entity
	DeleteEntity entity.Entity
	Reason       string
	Keyname      keyname.Keyname
}

// RegisterParams represents the register params
type RegisterParams struct {
	EntityMetaData entity.MetaData
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK          crypto.PrivateKey
	Client      applications.Client
	RoutePrefix string
}

var reg = createRegistry()

// SDKFunc represents the request SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Request
	Register             func(params RegisterParams)
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateSDKService     func(params CreateSDKServiceParams) Service
}{
	Create: func(params CreateParams) Request {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		if params.SaveEntity != nil {
			out := createRequestWithSaveEntity(params.ID, params.FromUser, params.SaveEntity, params.Reason, params.Keyname)
			return out
		}

		if params.DeleteEntity != nil {
			out := createRequestWithDeleteEntity(params.ID, params.FromUser, params.DeleteEntity, params.Reason, params.Keyname)
			return out
		}

		panic(errors.New("there is no Save or Delete entity in the request instance"))
	},
	Register: func(params RegisterParams) {
		regErr := reg.register(params.EntityMetaData)
		if regErr != nil {
			str := fmt.Sprintf("there was an error while registering an entity for a request: %s", regErr.Error())
			panic(str)
		}
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData(reg)
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(reg),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if req, ok := ins.(Request); ok {
					out, outErr := createStorableRequest(req)
					if outErr != nil {
						return nil, outErr
					}

					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if _, ok := ins.(Request); ok {
					return []string{
						retrieveAllRequestsKeyname(),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Request instance")
			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
				if req, ok := ins.(Request); ok {
					// metadata:
					metaData := createMetaData(reg)

					// create the repository and service:
					repository := entity.SDKFunc.CreateRepository(ds)
					keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
						EntityRepository: repository,
					})

					// make sure the request does not exists:
					_, retReqErr := repository.RetrieveByID(metaData, req.ID())
					if retReqErr == nil {
						str := fmt.Sprintf("the Request (ID: %s) already exists", req.ID().String())
						return errors.New(str)
					}

					// if the keyname does not exists, return an error:
					_, retKeynameErr := keynameRepository.RetrieveByName(req.Keyname().Name())
					if retKeynameErr != nil {
						return retKeynameErr
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.RoutePrefix)
		return out
	},
}
