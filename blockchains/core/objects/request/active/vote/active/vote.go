package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
)

type vote struct {
	UUID *uuid.UUID     `json:"id"`
	Vot  core_vote.Vote `json:"vote"`
	Pwr  int            `json:"power"`
}

func createVote(id *uuid.UUID, vot core_vote.Vote, power int) (Vote, error) {
	out := vote{
		UUID: id,
		Vot:  vot,
		Pwr:  power,
	}

	return &out, nil
}

func createVoteFromNormalized(normalized *normalizedVote) (Vote, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	votIns, votInsErr := core_vote.SDKFunc.CreateMetaData().Denormalize()(normalized.Vote)
	if votInsErr != nil {
		return nil, votInsErr
	}

	if vot, ok := votIns.(core_vote.Vote); ok {
		return createVote(&id, vot, normalized.Power)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid active Vote instance", votIns.ID().String())
	return nil, errors.New(str)
}

func createVoteFromStorable(rep entity.Repository, storable *storableVote) (Vote, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	voteID, voteIDErr := uuid.FromString(storable.VoteID)
	if voteIDErr != nil {
		return nil, voteIDErr
	}

	voteIns, voteInsErr := rep.RetrieveByID(core_vote.SDKFunc.CreateMetaData(), &voteID)
	if voteInsErr != nil {
		return nil, voteInsErr
	}

	if vot, ok := voteIns.(core_vote.Vote); ok {
		return createVote(&id, vot, storable.Power)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid active Vote instance", voteIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *vote) ID() *uuid.UUID {
	return obj.UUID
}

// Vote returns the vote
func (obj *vote) Vote() core_vote.Vote {
	return obj.Vot
}

// Power returns the power
func (obj *vote) Power() int {
	return obj.Pwr
}
