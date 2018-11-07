package vote

import (
	"errors"
	"fmt"

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

// CreateMetaDataParams represents the CreateMetaData params
type CreateMetaDataParams struct {
	Met entity.MetaData
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	Met entity.MetaData
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	CreateMetaData       func(params CreateMetaDataParams) entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
	CreateMetaData: func(params CreateMetaDataParams) entity.MetaData {
		return createMetaData(params.Met)
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(params.Met),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if vote, ok := ins.(Vote); ok {
					out := createStorableVote(vote)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Vote instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if vote, ok := ins.(Vote); ok {
					base := retrieveAllVotesKeyname()
					return []string{
						base,
						retrieveVotesByRequestIDKeyname(vote.Request().ID()),
						fmt.Sprintf("%s:by_voter_id:%s", base, vote.Voter().ID().String()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Vote instance")
			},
		})
	},
}
