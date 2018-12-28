package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Request represents an entity request
type Request interface {
	ID() *uuid.UUID
	From() user.User
	New() entity.Entity
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
	ID        *uuid.UUID
	FromUser  user.User
	NewEntity entity.Entity
	Reason    string
	Keyname   keyname.Keyname
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

// Data represents human-redable data
type Data struct {
	ID      string
	From    *user.Data
	New     string
	Reason  string
	Keyname *keyname.Data
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Requests    []*Data
}

var reg = createRegistry()

// SDKFunc represents the request SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Request
	Register             func(params RegisterParams)
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateSDKService     func(params CreateSDKServiceParams) Service
	ToData               func(req Request) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Request {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// get the keyname:
		out := createRequest(params.ID, params.FromUser, params.NewEntity, params.Reason, params.Keyname)
		return out
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
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				if req, ok := ins.(Request); ok {
					// metadata:
					metaData := createMetaData(reg)
					keynameRepresentation := keyname.SDKFunc.CreateRepresentation()

					// create the repository and service:
					repository := entity.SDKFunc.CreateRepository(ds)
					service := entity.SDKFunc.CreateService(ds)
					keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
						EntityRepository: repository,
					})

					// make sure the request does not exists:
					_, retReqErr := repository.RetrieveByID(metaData, req.ID())
					if retReqErr == nil {
						str := fmt.Sprintf("the Request (ID: %s) already exists", req.ID().String())
						return errors.New(str)
					}

					// if the keyname does not exists, create it:
					_, retKeynameErr := keynameRepository.RetrieveByName(req.Keyname().Name())
					if retKeynameErr != nil {
						saveErr := service.Save(req.Keyname(), keynameRepresentation)
						if saveErr != nil {
							return saveErr
						}
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
	ToData: func(req Request) *Data {
		out, outErr := toData(req)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
