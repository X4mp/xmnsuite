package task

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
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

// RetrieveSetByMilestone retrieves a Task set by milestone
func (app *repository) RetrieveSetByMilestone(mils milestone.Milestone, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllTaskKeyname(),
		retrieveTaskByMilestoneKeyname(mils),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}

// RetrieveSetByCreatedByUser retrieves a Task set by createdBy user
func (app *repository) RetrieveSetByCreatedByUser(crBy user.User, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllTaskKeyname(),
		retrieveTaskByCreatedByUserKeyname(crBy),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}
