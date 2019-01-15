package feature

type storableFeature struct {
	ID              string `json:"id"`
	ProjectID       string `json:"project_id"`
	Title           string `json:"title"`
	Details         string `json:"details"`
	CreatedByUserID string `json:"created_by_user_id"`
}

func createStorableFeature(ins Feature) *storableFeature {
	out := storableFeature{
		ID:              ins.ID().String(),
		ProjectID:       ins.Project().ID().String(),
		Title:           ins.Title(),
		Details:         ins.Details(),
		CreatedByUserID: ins.CreatedBy().ID().String(),
	}

	return &out
}
