package task

import "time"

type storableTask struct {
	ID              string    `json:"id"`
	MilestoneID     string    `json:"milestone_id"`
	CreatedByUserID string    `json:"created_by_user_id"`
	Title           string    `json:"title"`
	Details         string    `json:"details"`
	Deadline        time.Time `json:"deadline"`
	Reward          int       `json:"reward"`
	PledgeNeeded    int       `json:"pledge_needed"`
}

func createStorableTask(ins Task) *storableTask {
	out := storableTask{
		ID:              ins.ID().String(),
		MilestoneID:     ins.Milestone().ID().String(),
		CreatedByUserID: ins.CreatedBy().ID().String(),
		Title:           ins.Title(),
		Details:         ins.Details(),
		Deadline:        ins.Deadline(),
		Reward:          ins.Reward(),
		PledgeNeeded:    ins.PledgeNeeded(),
	}

	return &out
}
