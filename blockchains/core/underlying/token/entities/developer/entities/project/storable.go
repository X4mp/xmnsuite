package project

type storableProject struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createStorableProject(ins Project) *storableProject {
	out := storableProject{
		ID:          ins.ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
	}

	return &out
}
