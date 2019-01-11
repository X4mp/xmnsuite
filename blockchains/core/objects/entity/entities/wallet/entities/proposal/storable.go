package proposal

type storableProposal struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Details             string `json:"details"`
	CategoryID          string `json:"category_id"`
	ManagerPledgeNeeded int    `json:"manager_pledge_needed"`
	LinkerPledgeNeeded  int    `json:"linker_pledge_needed"`
}

func createStorableProposal(ins Proposal) *storableProposal {
	out := storableProposal{
		ID:                  ins.ID().String(),
		Title:               ins.Title(),
		Description:         ins.Description(),
		Details:             ins.Details(),
		CategoryID:          ins.Category().ID().String(),
		ManagerPledgeNeeded: ins.ManagerPledgeNeeded(),
		LinkerPledgeNeeded:  ins.LinkerPledgeNeeded(),
	}

	return &out
}
