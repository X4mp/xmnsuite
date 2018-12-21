package vote

import (
	uuid "github.com/satori/go.uuid"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// CalculateFn represents the vote calculation func.
// First bool = concensus is reached
// Second bool = the vote passed
type CalculateFn func(votes entity.PartialSet) (bool, bool, error)

// CreateRouteFn creates a route
type CreateRouteFn func(ins Vote, rep entity.Representation) (string, error)

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

// Repository represents the vote repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Vote, error)
	RetrieveByRequestVoter(req request.Request, voter user.User) (Vote, error)
	RetrieveSetByRequest(req request.Request, index int, amount int) (entity.PartialSet, error)
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
	CalculateVoteFn CalculateFn
	DS              datastore.DataStore
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK              crypto.PrivateKey
	Client          applications.Client
	CreateRouteFunc CreateRouteFn
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
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
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
	CreateService: func(params CreateServiceParams) Service {
		voteRepresentation := createRepresentation()
		requestRepresentation := request.SDKFunc.CreateRepresentation()
		entityRepository := entity.SDKFunc.CreateRepository(params.DS)
		entityService := entity.SDKFunc.CreateService(params.DS)
		out := createVoteService(params.CalculateVoteFn, entityRepository, entityService, voteRepresentation, requestRepresentation)
		return out
	},
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.CreateRouteFunc)
		return out
	},
}
