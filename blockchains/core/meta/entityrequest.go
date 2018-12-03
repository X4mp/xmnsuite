package meta

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/request/vote"
	"github.com/xmnservices/xmnsuite/datastore"
)

type entityRequest struct {
	ent                 entity.Representation
	mp                  map[string]entity.Representation
	createVoteServiceFn CreateVoteServiceFn
}

func createEntityRequest(ent entity.Representation, mp map[string]entity.Representation, createVoteServiceFn CreateVoteServiceFn) EntityRequest {
	out := entityRequest{
		ent:                 ent,
		mp:                  mp,
		createVoteServiceFn: createVoteServiceFn,
	}

	return &out
}

// RequestedBy returns the entity representation
func (obj *entityRequest) RequestedBy() entity.Representation {
	return obj.ent
}

// Map returns the map representation
func (obj *entityRequest) Map() map[string]entity.Representation {
	return obj.mp
}

// VoteService returns the vote service
func (obj *entityRequest) VoteService(store datastore.DataStore) vote.Service {
	return obj.createVoteServiceFn(store)
}
