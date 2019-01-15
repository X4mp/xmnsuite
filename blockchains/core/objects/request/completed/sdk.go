package completed

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Request represents a completed request
type Request interface {
	ID() *uuid.UUID
	Request() prev_request.Request
	ConcensusNeeded() int
	Approved() int
	Disapproved() int
	Neutral() int
}

// Normalized represents a normalized completed request
type Normalized interface {
}

// Repository represents the completed request repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Request, error)
	RetrieveByRequest(req prev_request.Request) (Request, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID              *uuid.UUID
	Request         prev_request.Request
	ConcensusNeeded int
	Approved        int
	Disapproved     int
	Neutral         int
}

// CreateRepositoryParams represents the create repository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the completed request SDK func
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

		out, outErr := createRequest(params.ID, params.Request, params.ConcensusNeeded, params.Approved, params.Disapproved, params.Neutral)
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
					out := createStorableRequest(req)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if req, ok := ins.(Request); ok {
					return []string{
						retrieveAllRequestsKeyname(),
						retrieveRequestByRequestKeyname(req.Request()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Request instance")
			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
				if req, ok := ins.(Request); ok {
					// metadata:
					requestMetaData := prev_request.SDKFunc.CreateMetaData()

					// create the repository and service:
					entityRepository := entity.SDKFunc.CreateRepository(ds)

					// make sure the request exists:
					_, retReqErr := entityRepository.RetrieveByID(requestMetaData, req.Request().ID())
					if retReqErr != nil {
						str := fmt.Sprintf("the completed request (ID: %s) contains a request instance (ID: %s) that does not exists", req.ID().String(), req.Request().ID().String())
						return errors.New(str)
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
}
