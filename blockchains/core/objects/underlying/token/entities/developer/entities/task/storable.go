package task

import (
	"time"
)

type storableTask struct {
	ID          string    `json:"id"`
	MilestoneID string    `json:"milestone_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_on"`
	DueOn       time.Time `json:"due_on"`
}

func createStorableTask(ins Task) *storableTask {
	out := storableTask{
		ID:          ins.ID().String(),
		MilestoneID: ins.Milestone().ID().String(),
		CreatorID:   ins.Creator().ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
		CreatedOn:   ins.CreatedOn(),
		DueOn:       ins.DueOn(),
	}

	return &out
}
