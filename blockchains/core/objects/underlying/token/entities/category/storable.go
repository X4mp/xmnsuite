package category

type storableCategory struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

func createStorableCategory(ins Category) *storableCategory {
	parentID := ""
	if ins.HasParent() {
		parentID = ins.Parent().ID().String()
	}

	out := storableCategory{
		ID:          ins.ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
		ParentID:    parentID,
	}

	return &out
}
