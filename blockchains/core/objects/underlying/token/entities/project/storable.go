package project

type storableProject struct {
	ID         string `json:"id"`
	ProposalID string `json:"proposal_id"`
}

func createStorableProject(ins Project) *storableProject {
	out := storableProject{
		ID:         ins.ID().String(),
		ProposalID: ins.Proposal().ID().String(),
	}

	return &out
}
