package project

type storableProject struct {
	ID         string `json:"id"`
	ProjectID string `json:"project_id"`
	OwnerID    string `json:"wallet_owner_id"`
	MgrID      string `json:"wallet_manager_id"`
	MgrShares  int    `json:"manager_shares"`
	LnkID      string `json:"wallet_linker_id"`
	LnkShares  int    `json:"linker_shares"`
	WrkShares  int    `json:"worker_shares"`
}

func createStorableProject(ins Project) *storableProject {
	out := storableProject{
		ID:         ins.ID().String(),
		ProjectID: ins.Project().ID().String(),
		OwnerID:    ins.Owner().ID().String(),
		MgrID:      ins.Manager().ID().String(),
		MgrShares:  ins.ManagerShares(),
		LnkID:      ins.Linker().ID().String(),
		LnkShares:  ins.LinkerShares(),
		WrkShares:  ins.WorkerShares(),
	}

	return &out
}
