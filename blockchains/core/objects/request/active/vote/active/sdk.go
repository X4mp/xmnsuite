package active

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Vote represents an active request vote
type Vote interface {
	ID() *uuid.UUID
	Vote() core_vote.Vote
	Power() int
}

// Normalized represents a normalized vote
type Normalized interface {
}

// Repository represents a vote repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Vote, error)
	RetrieveByVote(vot core_vote.Vote) (Vote, error)
	RetrieveByRequestVoter(voter user.User, req active_request.Request) (Vote, error)
	RetrieveSetByRequest(req active_request.Request, index int, amount int) (entity.PartialSet, error)
}

// Service represents the vote service
type Service interface {
	Save(ins Vote, rep entity.Representation) error
}

// CreateParams represents the create params
type CreateParams struct {
	ID    *uuid.UUID
	Vote  core_vote.Vote
	Power int
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	DS               datastore.DataStore
	EntityRepository entity.Repository
	EntityService    entity.Service
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
}{
	Create: func(params CreateParams) Vote {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// create the request:
		out, outErr := createVote(params.ID, params.Vote, params.Power)
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
		if params.EntityService == nil && params.EntityRepository == nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.DS)
			params.EntityService = entity.SDKFunc.CreateService(params.DS)
		}

		voteRepresentation := createRepresentation()
		requestRepresentation := active_request.SDKFunc.CreateRepresentation()
		out := createVoteService(params.EntityRepository, params.EntityService, voteRepresentation, requestRepresentation)
		return out
	},
}
