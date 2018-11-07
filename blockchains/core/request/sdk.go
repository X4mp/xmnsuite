package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

// Request represents an entity request
type Request interface {
	ID() *uuid.UUID
	From() user.User
	New() entity.Entity
}

// CreateParams represents the create params
type CreateParams struct {
	ID        *uuid.UUID
	FromUser  user.User
	NewEntity entity.Entity
}

// CreateMetaDataParams represents the CreateMetaData params
type CreateMetaDataParams struct {
	Met entity.MetaData
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	Met entity.MetaData
}

// SDKFunc represents the request SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Request
	CreateMetaData       func(params CreateMetaDataParams) entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
	Create: func(params CreateParams) Request {
		out := createRequest(params.ID, params.FromUser, params.NewEntity)
		return out
	},
	CreateMetaData: func(params CreateMetaDataParams) entity.MetaData {
		return createMetaData(params.Met)
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(params.Met),
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
}
