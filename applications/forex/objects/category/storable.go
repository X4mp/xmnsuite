package category

type storableCategory struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createStorableCategory(ins Category) *storableCategory {

	parentIDAsString := ""
	if ins.HasParent() {
		parentIDAsString = ins.Parent().ID().String()
	}

	out := storableCategory{
		ID:          ins.ID().String(),
		ParentID:    parentIDAsString,
		Name:        ins.Name(),
		Description: ins.Description(),
	}

	return &out
}
