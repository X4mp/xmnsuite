package completed

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

type task struct {
	UUID *uuid.UUID     `json:"id"`
	Tsk  mils_task.Task `json:"task"`
	Det  string         `json:"details"`
}

func createTask(id *uuid.UUID, tsk mils_task.Task, details string) (Task, error) {
	out := task{
		UUID: id,
		Tsk:  tsk,
		Det:  details,
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

	if tsk, ok := tskIns.(mils_task.Task); ok {
		return createTask(&id, tsk, normalized.Details)
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

	tskIns, tskInsErr := rep.RetrieveByID(mils_task.SDKFunc.CreateMetaData(), &taskID)
	if tskInsErr != nil {
		return nil, tskInsErr
	}

	if tsk, ok := tskIns.(mils_task.Task); ok {
		return createTask(&id, tsk, storable.Details)
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

// Details returns the details
func (obj *task) Details() string {
	return obj.Det
}
