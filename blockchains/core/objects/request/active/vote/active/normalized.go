package active

import (
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
)

type normalizedVote struct {
	ID    string               `json:"id"`
	Vote  core_vote.Normalized `json:"vote"`
	Power int                  `json:"power"`
}

func createNormalizedVote(ins Vote) (*normalizedVote, error) {
	vot, votErr := core_vote.SDKFunc.CreateMetaData().Normalize()(ins.Vote())
	if votErr != nil {
		return nil, votErr
	}

	out := normalizedVote{
		ID:    ins.ID().String(),
		Vote:  vot,
		Power: ins.Power(),
	}

	return &out, nil
}
