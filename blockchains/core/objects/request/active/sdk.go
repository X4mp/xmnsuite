package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	core_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Request represents an active request
type Request interface {
	ID() *uuid.UUID
	Request() core_request.Request
	ConcensusNeeded() int
}

// Normalized represents a normalized request
type Normalized interface {
}

// Repository represents a Request repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Request, error)
	RetrieveByRequest(req core_request.Request) (Request, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByFromUser(usr user.User, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByKeyname(kname keyname.Keyname, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID              *uuid.UUID
	Request         core_request.Request
	ConcensusNeeded int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the request SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Request
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Request {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// create the request:
		out, outErr := createRequest(params.ID, params.Request, params.ConcensusNeeded)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if req, ok := ins.(Request); ok {
					out := createStorable(req)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Request instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if req, ok := ins.(Request); ok {
					return []string{
						retrieveAllRequestsKeyname(),
						retrieveAllRequestsByRequestKeyname(req.Request()),
						retrieveAllRequestsFromUserKeyname(req.Request().From()),
						retrieveAllRequestsByKeynameKeyname(req.Request().Keyname()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid active Request instance")
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				if req, ok := ins.(Request); ok {
					// metadata:
					metaData := createMetaData()
					coreRequestRepresentation := core_request.SDKFunc.CreateRepresentation()

					// create the repository and service:
					entityRepository := entity.SDKFunc.CreateRepository(ds)
					service := entity.SDKFunc.CreateService(ds)

					// make sure the request does not exists:
					_, retReqErr := entityRepository.RetrieveByID(metaData, req.ID())
					if retReqErr == nil {
						str := fmt.Sprintf("the Request (ID: %s) already exists", req.ID().String())
						return errors.New(str)
					}

					// make sure the request does not exits, then save it:
					_, retPrevReqErr := entityRepository.RetrieveByID(coreRequestRepresentation.MetaData(), req.Request().ID())
					if retPrevReqErr == nil {
						str := fmt.Sprintf("the given Request (ID: %s) already exists: %s", req.Request().ID().String(), retPrevReqErr.Error())
						return errors.New(str)
					}

					// save the request:
					saveReqErr := service.Save(req.Request(), coreRequestRepresentation)
					if saveReqErr != nil {
						return saveReqErr
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Request instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		return createRepository(params.EntityRepository, metaData)
	},
}
