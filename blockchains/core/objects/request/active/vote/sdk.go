package vote

import (
	uuid "github.com/satori/go.uuid"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/crypto"
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
	Reason() string
	IsNeutral() bool
	IsApproved() bool
}

// Normalized represents a normalized Vote
type Normalized interface {
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
	Reason     string
	IsApproved bool
	IsNeutral  bool
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK          crypto.PrivateKey
	Client      applications.Client
	RoutePrefix string
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateSDKService     func(params CreateSDKServiceParams) Service
}{
	Create: func(params CreateParams) Vote {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createVote(params.ID, params.Request, params.Voter, params.Reason, params.IsNeutral, params.IsApproved)
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
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.RoutePrefix)
		return out
	},
}
