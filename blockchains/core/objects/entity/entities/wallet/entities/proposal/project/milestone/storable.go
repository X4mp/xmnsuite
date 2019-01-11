package milestone

type storableMilestone struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	FeatureID   string `json:"feature_id"`
	WalletID    string `json:"wallet_id"`
	Shares      int    `json:"shares"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func createStorableMilestone(ins Milestone) *storableMilestone {
	featID := ""
	if ins.HasFeature() {
		featID = ins.Feature().ID().String()
	}

	out := storableMilestone{
		ID:          ins.ID().String(),
		ProjectID:   ins.Project().ID().String(),
		FeatureID:   featID,
		WalletID:    ins.Wallet().ID().String(),
		Shares:      ins.Shares(),
		Title:       ins.Title(),
		Description: ins.Description(),
		Details:     ins.Details(),
	}

	return &out
}
