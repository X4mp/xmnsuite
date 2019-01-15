package task

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllTaskKeyname() string {
	return "tasks"
}

func retrieveTaskByMilestoneKeyname(mils milestone.Milestone) string {
	base := retrieveAllTaskKeyname()
	return fmt.Sprintf("%s:by_milestone_id:%s", base, mils.ID().String())
}

func retrieveTaskByCreatedByUserKeyname(createdBy user.User) string {
	base := retrieveAllTaskKeyname()
	return fmt.Sprintf("%s:by_created_by_user_id:%s", base, createdBy.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Task",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableTask); ok {
				return createTaskFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedTask)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createTaskFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Task instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if tsk, ok := ins.(Task); ok {
				out, outErr := createNormalizedTask(tsk)
				if outErr != nil {
					return nil, outErr
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

			return nil, errors.New("the given normalized instance cannot be converted to a Task instance")
		},
		EmptyStorable:   new(storableTask),
		EmptyNormalized: new(normalizedTask),
	})
}

func representation() entity.Representation {
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
			if tsk, ok := ins.(Task); ok {
				return []string{
					retrieveAllTaskKeyname(),
					retrieveTaskByMilestoneKeyname(tsk.Milestone()),
					retrieveTaskByCreatedByUserKeyname(tsk.CreatedBy()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if tsk, ok := ins.(Task); ok {
				// crate metadata and representation:
				milestoneMetaData := milestone.SDKFunc.CreateMetaData()
				userMetaData := user.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)

				// make sure the milestone exists:
				_, retProjErr := entityRepository.RetrieveByID(milestoneMetaData, tsk.Milestone().ID())
				if retProjErr != nil {
					str := fmt.Sprintf("the given task (ID: %s) contains a milestone (ID: %s) that does not exists", tsk.ID().String(), tsk.Milestone().ID().String())
					return errors.New(str)
				}

				// make sure the user exists:
				_, retWalErr := entityRepository.RetrieveByID(userMetaData, tsk.CreatedBy().ID())
				if retWalErr != nil {
					str := fmt.Sprintf("the given task (ID: %s) contains a createdBy user (ID: %s) that does not exists", tsk.ID().String(), tsk.CreatedBy().ID().String())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Task instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
