package proposal

import "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"

type normalizedProposal struct {
	ID                  string              `json:"id"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	Details             string              `json:"details"`
	Category            category.Normalized `json:"category"`
	ManagerPledgeNeeded int                 `json:"manager_pledge_needed"`
	LinkerPledgeNeeded  int                 `json:"linker_pledge_needed"`
}

func createNormalizedProposal(ins Proposal) (*normalizedProposal, error) {
	cat, catErr := category.SDKFunc.CreateMetaData().Normalize()(ins.Category())
	if catErr != nil {
		return nil, catErr
	}

	out := normalizedProposal{
		ID:                  ins.ID().String(),
		Title:               ins.Title(),
		Description:         ins.Description(),
		Details:             ins.Details(),
		Category:            cat,
		ManagerPledgeNeeded: ins.ManagerPledgeNeeded(),
		LinkerPledgeNeeded:  ins.LinkerPledgeNeeded(),
	}

	return &out, nil
}
