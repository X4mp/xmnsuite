package completed

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
	pledge_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllTaskKeyname() string {
	return "completedtasks"
}

func retrieveTaskByTaskKeyname(tsk mils_task.Task) string {
	base := retrieveAllTaskKeyname()
	return fmt.Sprintf("%s:by_task_id:%s", base, tsk.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "CompletedTask",
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

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid completed Task instance", ins.ID().String())
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

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid completed Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if tsk, ok := ins.(Task); ok {
				return []string{
					retrieveAllTaskKeyname(),
					retrieveTaskByTaskKeyname(tsk.Task()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid completed Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if tsk, ok := ins.(Task); ok {
				// crate metadata and representation:
				pledgeTaskRepresentation := pledge_task.SDKFunc.CreateRepresentation()
				depositRepresentation := deposit.SDKFunc.CreateRepresentation()
				taskMetaData := mils_task.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				entityService := entity.SDKFunc.CreateService(ds)
				pledgeTaskRepository := pledge_task.SDKFunc.CreateRepository(pledge_task.CreateRepositoryParams{
					EntityRepository: entityRepository,
				})

				// make sure the task exists:
				_, retTskErr := entityRepository.RetrieveByID(taskMetaData, tsk.Task().ID())
				if retTskErr != nil {
					str := fmt.Sprintf("the given completed task (ID: %s) contains a task (ID: %s) that does not exists", tsk.ID().String(), tsk.Task().ID().String())
					return errors.New(str)
				}

				//retrieve the pledge task:
				pldgeTask, pldgeTaskErr := pledgeTaskRepository.RetrieveByTask(tsk.Task())
				if pldgeTaskErr != nil {
					str := fmt.Sprintf("there was an error while retrieving the pledge task from milestone tasl (ID: %s): %s", tsk.Task().ID().String(), pldgeTaskErr.Error())
					return errors.New(str)
				}

				// create the deposit:
				frm := pldgeTask.Pledge().From()
				dep := deposit.SDKFunc.Create(deposit.CreateParams{
					To:     frm.From(),
					Token:  frm.Token(),
					Amount: pldgeTask.Task().Reward(),
				})

				// save the deposit:
				saveDepErr := entityService.Save(dep, depositRepresentation)
				if saveDepErr != nil {
					str := fmt.Sprintf("there was an error while saving the deposit (ID: %s): %s", dep.ID().String(), saveDepErr.Error())
					return errors.New(str)
				}

				// delete the pledge task:
				delPledgeTaskErr := entityService.Delete(pldgeTask, pledgeTaskRepresentation)
				if delPledgeTaskErr != nil {
					str := fmt.Sprintf("there was an error while deleting the pledge task (ID: %s): %s", pldgeTask.ID().String(), delPledgeTaskErr.Error())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid completed Task instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
