package vote

import (
	uuid "github.com/satori/go.uuid"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

// Vote represents a request vote
type Vote interface {
	ID() *uuid.UUID
	Request() request.Request
	Voter() user.User
	IsApproved() bool
}

// Service represents the vote service
type Service interface {
	Save(ins Vote) error
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Request    request.Request
	Voter      user.User
	IsApproved bool
}

// CreateMetaDataParams represents the CreateMetaData params
type CreateMetaDataParams struct {
	RequestMetaData entity.MetaData
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	RequestMetaData entity.MetaData
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	EntityRepository        entity.Repository
	EntityService           entity.Service
	RequestRepresentation   entity.Representation
	NewEntityRepresentation entity.Representation
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func(params CreateMetaDataParams) entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
	CreateService        func(params CreateServiceParams) Service
}{
	Create: func(params CreateParams) Vote {
		out, outErr := createVote(params.ID, params.Request, params.Voter, params.IsApproved)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func(params CreateMetaDataParams) entity.MetaData {
		return createMetaData(params.RequestMetaData)
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
		return createRepresentation(params.RequestMetaData)
	},
	CreateService: func(params CreateServiceParams) Service {
		voteRepresentation := createRepresentation(params.RequestRepresentation.MetaData())
		out := createVoteService(params.EntityRepository, params.EntityService, voteRepresentation, params.RequestRepresentation, params.NewEntityRepresentation)
		return out
	},
}
