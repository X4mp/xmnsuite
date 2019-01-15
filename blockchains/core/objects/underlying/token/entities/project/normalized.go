package project

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
)

type normalizedProject struct {
	ID       string              `json:"id"`
	Proposal proposal.Normalized `json:"proposal"`
}

func createNormalizedProject(ins Project) (*normalizedProject, error) {
	prop, propErr := proposal.SDKFunc.CreateMetaData().Normalize()(ins.Proposal())
	if propErr != nil {
		return nil, propErr
	}

	out := normalizedProject{
		ID:       ins.ID().String(),
		Proposal: prop,
	}

	return &out, nil
}
