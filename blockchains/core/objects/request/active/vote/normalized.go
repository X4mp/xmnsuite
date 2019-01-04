package vote

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
)

type normalizedVote struct {
	ID         string             `json:"id"`
	Request    request.Normalized `json:"request"`
	Voter      user.Normalized    `json:"voter"`
	Reason     string             `json:"reason"`
	IsNeutral  bool               `json:"is_neutral"`
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
		Reason:     ins.Reason(),
		IsNeutral:  ins.IsNeutral(),
		IsApproved: ins.IsApproved(),
	}

	return &out, nil
}
