package completed

import (
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type normalizedRequest struct {
	ID              string                  `json:"id"`
	Request         prev_request.Normalized `json:"request"`
	ConcensusNeeded int                     `json:"concensus_needed"`
	Approved        int                     `json:"approved"`
	Disapproved     int                     `json:"disapproved"`
	Neutral         int                     `json:"neutral"`
}

func createNormalizedRequest(ins Request) (*normalizedRequest, error) {
	req, reqErr := prev_request.SDKFunc.CreateMetaData().Normalize()(ins.Request())
	if reqErr != nil {
		return nil, reqErr
	}

	out := normalizedRequest{
		ID:              ins.ID().String(),
		Request:         req,
		ConcensusNeeded: ins.ConcensusNeeded(),
		Approved:        ins.Approved(),
		Disapproved:     ins.Disapproved(),
		Neutral:         ins.Neutral(),
	}

	return &out, nil
}
