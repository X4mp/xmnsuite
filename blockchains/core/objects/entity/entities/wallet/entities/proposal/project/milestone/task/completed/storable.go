package completed

type storableTask struct {
	ID      string `json:"id"`
	TaskID  string `json:"task_id"`
	Details string `json:"details"`
}

func createStorableTask(ins Task) *storableTask {
	out := storableTask{
		ID:      ins.ID().String(),
		TaskID:  ins.Task().ID().String(),
		Details: ins.Details(),
	}

	return &out
}
