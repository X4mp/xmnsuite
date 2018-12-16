package vote

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type normalizedVote struct {
	ID         string             `json:"id"`
	Request    request.Normalized `json:"request"`
	Voter      user.Normalized    `json:"voter"`
	IsApproved bool               `json:"is_approved"`
}

func createNormalizedVote(ins Vote) (*normalizedVote, error) {
	req, reqErr := request.SDKFunc.CreateMetaData().Normalize()(ins.Request())
	if reqErr != nil {
		return nil, reqErr
	}

	voter, voterErr := user.SDKFunc.CreateMetaData().Normalize()(ins.Voter())
	if voterErr != nil {
		return nil, voterErr
	}

	out := normalizedVote{
		ID:         ins.ID().String(),
		Request:    req,
		Voter:      voter,
		IsApproved: ins.IsApproved(),
	}

	return &out, nil
}
