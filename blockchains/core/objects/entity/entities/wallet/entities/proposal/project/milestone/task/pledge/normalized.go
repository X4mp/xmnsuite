package pledge

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

type normalizedTask struct {
	ID     string               `json:"id"`
	Task   mils_task.Normalized `json:"task"`
	Pledge pledge.Normalized    `json:"pledge"`
}

func createNormalizedTask(ins Task) (*normalizedTask, error) {
	tsk, tskErr := mils_task.SDKFunc.CreateMetaData().Normalize()(ins.Task())
	if tskErr != nil {
		return nil, tskErr
	}

	pldge, pldgeErr := pledge.SDKFunc.CreateMetaData().Normalize()(ins.Pledge())
	if pldgeErr != nil {
		return nil, pldgeErr
	}

	out := normalizedTask{
		ID:     ins.ID().String(),
		Task:   tsk,
		Pledge: pldge,
	}

	return &out, nil
}
