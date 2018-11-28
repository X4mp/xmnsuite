package milestone

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
)

func retrieveAllMilestonesKeyname() string {
	return "milestones"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Milestone",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableMilestone); ok {
				// metadata:
				projectMetaData := project.SDKFunc.CreateMetaData()

				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				projectID, projectIDErr := uuid.FromString(storable.ProjectID)
				if projectIDErr != nil {
					return nil, projectIDErr
				}

				projIns, projInsErr := rep.RetrieveByID(projectMetaData, &projectID)
				if projInsErr != nil {
					return nil, projInsErr
				}

				if proj, ok := projIns.(project.Project); ok {
					out := createMilestone(&id, proj, storable.Title, storable.Description, storable.CreatedOn, storable.DueOn)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance and thererfore the given data cannot be transformed to a Milestone instance", projIns.ID().String())
				return nil, errors.New(str)

			}

			ptr := new(normalizedMilestone)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createMilestoneFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if milestone, ok := ins.(Milestone); ok {
				out, outErr := createNormalizedMilestone(milestone)
				return out, outErr
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedMilestone); ok {
				return createMilestoneFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Milestone instance")
		},
		EmptyStorable:   new(storableMilestone),
		EmptyNormalized: new(normalizedMilestone),
	})
}
