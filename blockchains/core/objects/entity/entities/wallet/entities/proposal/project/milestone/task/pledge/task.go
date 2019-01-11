package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

type task struct {
	UUID  *uuid.UUID     `json:"id"`
	Tsk   mils_task.Task `json:"task"`
	Pldge pledge.Pledge  `json:"pledge"`
}

func createTask(id *uuid.UUID, tsk mils_task.Task, pldge pledge.Pledge) (Task, error) {
	out := task{
		UUID:  id,
		Tsk:   tsk,
		Pldge: pldge,
	}

	return &out, nil
}

func createTaskFromNormalized(normalized *normalizedTask) (Task, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	tskIns, tskInsErr := mils_task.SDKFunc.CreateMetaData().Denormalize()(normalized.Task)
	if tskInsErr != nil {
		return nil, tskInsErr
	}

	pldgeIns, pldgeInsErr := pledge.SDKFunc.CreateMetaData().Denormalize()(normalized.Pledge)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	if tsk, ok := tskIns.(mils_task.Task); ok {
		if pldge, ok := pldgeIns.(pledge.Pledge); ok {
			return createTask(&id, tsk, pldge)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid pledge Task instance", tskIns.ID().String())
	return nil, errors.New(str)
}

func createTaskFromStorable(storable *storableTask, rep entity.Repository) (Task, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	taskID, taskIDErr := uuid.FromString(storable.TaskID)
	if taskIDErr != nil {
		return nil, taskIDErr
	}

	pldgeID, pldgeIDErr := uuid.FromString(storable.PledgeID)
	if pldgeIDErr != nil {
		return nil, pldgeIDErr
	}

	tskIns, tskInsErr := rep.RetrieveByID(mils_task.SDKFunc.CreateMetaData(), &taskID)
	if tskInsErr != nil {
		return nil, tskInsErr
	}

	pldgeIns, pldgeInsErr := rep.RetrieveByID(pledge.SDKFunc.CreateMetaData(), &pldgeID)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	if tsk, ok := tskIns.(mils_task.Task); ok {
		if pldge, ok := pldgeIns.(pledge.Pledge); ok {
			return createTask(&id, tsk, pldge)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid pledge Task instance", tskIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *task) ID() *uuid.UUID {
	return obj.UUID
}

// Task returns the task
func (obj *task) Task() mils_task.Task {
	return obj.Tsk
}

// Pledge returns the pledge
func (obj *task) Pledge() pledge.Pledge {
	return obj.Pldge
}
