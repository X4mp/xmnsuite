package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

type repository struct {
	metaData         entity.MetaData
	entityRepository entity.Repository
}

func createRepository(metaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		metaData:         metaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves a Task by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Task, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if tsk, ok := ins.(Task); ok {
		return tsk, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Task instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByTask retrieves a Task by Task
func (app *repository) RetrieveByTask(tsk mils_task.Task) (Task, error) {
	keynames := []string{
		retrieveAllTaskKeyname(),
		retrieveTaskByTaskKeyname(tsk),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if tsk, ok := ins.(Task); ok {
		return tsk, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Task instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByPledge retrieves a Task by Pledge
func (app *repository) RetrieveByPledge(pldge pledge.Pledge) (Task, error) {
	keynames := []string{
		retrieveAllTaskKeyname(),
		retrieveTaskByPledgeKeyname(pldge),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if tsk, ok := ins.(Task); ok {
		return tsk, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Task instance", ins.ID().String())
	return nil, errors.New(str)
}
