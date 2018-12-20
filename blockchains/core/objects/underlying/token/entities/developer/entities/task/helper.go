package task

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/milestone"
)

func retrieveAllTasksKeyname() string {
	return "tasks"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Task",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableTask); ok {
				// metadata:
				milestoneMetadata := milestone.SDKFunc.CreateMetaData()
				devMetaData := developer.SDKFunc.CreateMetaData()

				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				milestoneID, milestoneIDErr := uuid.FromString(storable.MilestoneID)
				if milestoneIDErr != nil {
					return nil, milestoneIDErr
				}

				milestoneIns, milestoneInsErr := rep.RetrieveByID(milestoneMetadata, &milestoneID)
				if milestoneInsErr != nil {
					return nil, milestoneInsErr
				}

				creatorID, creatorIDErr := uuid.FromString(storable.CreatorID)
				if creatorIDErr != nil {
					return nil, creatorIDErr
				}

				creatorIns, creatorInsErr := rep.RetrieveByID(devMetaData, &creatorID)
				if creatorInsErr != nil {
					return nil, creatorInsErr
				}

				if mils, ok := milestoneIns.(milestone.Milestone); ok {
					if creator, ok := creatorIns.(developer.Developer); ok {
						out := createTask(&id, mils, creator, storable.Title, storable.Description, storable.CreatedOn, storable.DueOn)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Developer instance and thererfore the given data cannot be transformed to a Task instance", creatorIns.ID().String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance and thererfore the given data cannot be transformed to a Task instance", milestoneIns.ID().String())
				return nil, errors.New(str)

			}

			ptr := new(normalizedTask)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createTaskFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if tsk, ok := ins.(Task); ok {
				out, outErr := createNormalizedTask(tsk)
				if outErr != nil {
					panic(outErr)
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedTask); ok {
				return createTaskFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Task instance")
		},
		EmptyStorable:   new(storableTask),
		EmptyNormalized: new(normalizedTask),
	})
}
