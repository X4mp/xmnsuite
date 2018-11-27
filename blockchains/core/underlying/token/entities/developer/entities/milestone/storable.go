package milestone

import "time"

type storableMilestone struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_on"`
	DueOn       time.Time `json:"due_on"`
}

func createStorableMilestone(ins Milestone) *storableMilestone {
	out := storableMilestone{
		ID:          ins.ID().String(),
		ProjectID:   ins.Project().ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
		CreatedOn:   ins.CreatedOn(),
		DueOn:       ins.DueOn(),
	}

	return &out
}
