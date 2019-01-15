package pledge

type storableTask struct {
	ID       string `json:"id"`
	TaskID   string `json:"task_id"`
	PledgeID string `json:"pledge_id"`
}

func createStorableTask(ins Task) *storableTask {
	out := storableTask{
		ID:       ins.ID().String(),
		TaskID:   ins.Task().ID().String(),
		PledgeID: ins.Pledge().ID().String(),
	}

	return &out
}
