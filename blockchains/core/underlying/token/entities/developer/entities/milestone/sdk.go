package milestone

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Milestone represents a milestone
type Milestone interface {
	ID() *uuid.UUID
	Project() project.Project
	Title() string
	Description() string
	CreatedOn() time.Time
	DueOn() time.Time
}

// Normalized represents a normalized milestone
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID          *uuid.UUID
	Project     project.Project
	Title       string
	Description string
	CreatedOn   time.Time
	DueOn       time.Time
}

// SDKFunc represents the Developer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Milestone
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Milestone {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createMilestone(params.ID, params.Project, params.Title, params.Description, params.CreatedOn, params.DueOn)
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
				if milestone, ok := ins.(Milestone); ok {
					out := createStorableMilestone(milestone)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if _, ok := ins.(Milestone); ok {
					return []string{
						retrieveAllMilestonesKeyname(),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the metadata:
				metaData := createMetaData()
				projectMetaData := project.SDKFunc.CreateMetaData()

				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				if milestone, ok := ins.(Milestone); ok {
					// if the milestone already exists, return an error:
					_, retMilestoneErr := repository.RetrieveByID(metaData, milestone.ID())
					if retMilestoneErr == nil {
						str := fmt.Sprintf("the Milestone (ID: %s) already exists", milestone.ID().String())
						return errors.New(str)
					}

					// if the project does not exists, return an error:
					proj := milestone.Project()
					_, retProjErr := repository.RetrieveByID(projectMetaData, proj.ID())
					if retProjErr != nil {
						str := fmt.Sprintf("the Project (ID: %s) in the Milestone instance (ID: %s) does not exists", proj.ID().String(), milestone.ID().String())
						return errors.New(str)
					}

					// everything is alright:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
