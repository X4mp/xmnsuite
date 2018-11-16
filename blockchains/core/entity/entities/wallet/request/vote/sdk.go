package vote

import (
	uuid "github.com/satori/go.uuid"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Vote represents a request vote
type Vote interface {
	ID() *uuid.UUID
	Request() request.Request
	Voter() user.User
	IsApproved() bool
}

// NormalizedVote represents a normalized Vote
type NormalizedVote interface {
}

// Service represents the vote service
type Service interface {
	Save(ins Vote, rep entity.Representation) error
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Request    request.Request
	Voter      user.User
	IsApproved bool
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	EntityRepository entity.Repository
	EntityService    entity.Service
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK     crypto.PrivateKey
	Client applications.Client
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateService        func(params CreateServiceParams) Service
	CreateSDKService     func(params CreateSDKServiceParams) Service
}{
	Create: func(params CreateParams) Vote {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createVote(params.ID, params.Request, params.Voter, params.IsApproved)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateService: func(params CreateServiceParams) Service {
		voteRepresentation := createRepresentation()
		requestRepresentation := request.SDKFunc.CreateRepresentation()
		out := createVoteService(params.EntityRepository, params.EntityService, voteRepresentation, requestRepresentation)
		return out
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client)
		return out
	},
}
