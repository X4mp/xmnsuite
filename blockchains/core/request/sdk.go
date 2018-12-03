package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Request represents an entity request
type Request interface {
	ID() *uuid.UUID
	From() user.User
	New() entity.Entity
	NewName() string
}

// Normalized represents a normalized request
type Normalized interface {
}

// Service represents an entity service
type Service interface {
	Save(req Request, entityRep entity.Representation) error
}

// Repository represents a Request repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Request, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByFromUser(usr user.User, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID             *uuid.UUID
	FromUser       user.User
	NewEntity      entity.Entity
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
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateSDKService     func(params CreateSDKServiceParams) Service
}{
	Create: func(params CreateParams) Request {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// get the keyname:
		name := params.EntityMetaData.Keyname()

		// register:
		regErr := reg.register(params.EntityMetaData)
		if regErr != nil {
			str := fmt.Sprintf("there was an error while registering an entity for a request: %s", regErr.Error())
			panic(str)
		}

		out := createRequest(params.ID, params.FromUser, params.NewEntity, name)
		return out
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
				if req, ok := ins.(Request); ok {
					base := retrieveAllRequestsKeyname()
					return []string{
						base,
						fmt.Sprintf("%s:by_from_id:%s", base, req.From().ID().String()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Request instance")
			},
		})
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.RoutePrefix)
		return out
	},
}
