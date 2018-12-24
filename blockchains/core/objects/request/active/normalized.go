package active

import (
	core_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type normalizedRequest struct {
	ID              string                  `json:"id"`
	Request         core_request.Normalized `json:"request"`
	ConcensusNeeded int                     `json:"concensus_needed"`
}

func createNormalizedRequest(ins Request) (*normalizedRequest, error) {
	req, reqErr := core_request.SDKFunc.CreateMetaData().Normalize()(ins.Request())
	if reqErr != nil {
		return nil, reqErr
	}

	out := normalizedRequest{
		ID:              ins.ID().String(),
		Request:         req,
		ConcensusNeeded: ins.ConcensusNeeded(),
	}

	return &out, nil
}
