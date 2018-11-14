package category

type storableCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createStorableCategory(ins Category) *storableCategory {
	out := storableCategory{
		ID:          ins.ID().String(),
		Name:        ins.Name(),
		Description: ins.Description(),
	}

	return &out
}
