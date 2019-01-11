package completed

import (
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

type normalizedTask struct {
	ID      string               `json:"id"`
	Task    mils_task.Normalized `json:"task"`
	Details string               `json:"details"`
}

func createNormalizedTask(ins Task) (*normalizedTask, error) {
	tsk, tskErr := mils_task.SDKFunc.CreateMetaData().Normalize()(ins.Task())
	if tskErr != nil {
		return nil, tskErr
	}

	out := normalizedTask{
		ID:      ins.ID().String(),
		Task:    tsk,
		Details: ins.Details(),
	}

	return &out, nil
}
