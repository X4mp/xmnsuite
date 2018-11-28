package task

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Task represents a task
type Task interface {
	ID() *uuid.UUID
	Milestone() milestone.Milestone
	Creator() developer.Developer
	Title() string
	Description() string
	CreatedOn() time.Time
	DueOn() time.Time
}

// Normalized represents a normalized task
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID          *uuid.UUID
	Milestone   milestone.Milestone
	Creator     developer.Developer
	Title       string
	Description string
	CreatedOn   time.Time
	DueOn       time.Time
}

// SDKFunc represents the Developer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Task
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Task {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createTask(params.ID, params.Milestone, params.Creator, params.Title, params.Description, params.CreatedOn, params.DueOn)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if tsk, ok := ins.(Task); ok {
					out := createStorableTask(tsk)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Task instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if _, ok := ins.(Task); ok {
					return []string{
						retrieveAllTasksKeyname(),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Task instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the metadata:
				metaData := createMetaData()
				milestoneMetaData := milestone.SDKFunc.CreateMetaData()

				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				if tsk, ok := ins.(Task); ok {
					// if the task already exists, return an error:
					_, retTaskErr := repository.RetrieveByID(metaData, tsk.ID())
					if retTaskErr == nil {
						str := fmt.Sprintf("the Task (ID: %s) already exists", tsk.ID().String())
						return errors.New(str)
					}

					// if the milestone does not exists, return an error:
					mstone := tsk.Milestone()
					_, retMilestoneErr := repository.RetrieveByID(milestoneMetaData, mstone.ID())
					if retMilestoneErr != nil {
						str := fmt.Sprintf("the Milestone (ID: %s) in the Task instance (ID: %s) does not exists", mstone.ID().String(), tsk.ID().String())
						return errors.New(str)
					}

					// everything is alright:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Task instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
