package pledge

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllTaskKeyname() string {
	return "pledgetasks"
}

func retrieveTaskByTaskKeyname(tsk mils_task.Task) string {
	base := retrieveAllTaskKeyname()
	return fmt.Sprintf("%s:by_task_id:%s", base, tsk.ID().String())
}

func retrieveTaskByPledgeKeyname(pldge pledge.Pledge) string {
	base := retrieveAllTaskKeyname()
	return fmt.Sprintf("%s:by_pledge_id:%s", base, pldge.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "PledgeTask",
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

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid pledge Task instance", ins.ID().String())
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

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid pledge Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if tsk, ok := ins.(Task); ok {
				return []string{
					retrieveAllTaskKeyname(),
					retrieveTaskByTaskKeyname(tsk.Task()),
					retrieveTaskByPledgeKeyname(tsk.Pledge()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid pledge Task instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if tsk, ok := ins.(Task); ok {
				// crate metadata and representation:
				metaData := createMetaData()
				taskMetaData := mils_task.SDKFunc.CreateMetaData()
				pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()
				pledgeMetaData := pledgeRepresentation.MetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				entityService := entity.SDKFunc.CreateService(ds)
				repository := createRepository(metaData, entityRepository)

				// make sure the task exists:
				_, retTskErr := entityRepository.RetrieveByID(taskMetaData, tsk.Task().ID())
				if retTskErr != nil {
					str := fmt.Sprintf("the given pledge task (ID: %s) contains a task (ID: %s) that does not exists", tsk.ID().String(), tsk.Task().ID().String())
					return errors.New(str)
				}

				// make sure the pledge task is not already assigned to another task:
				_, retPledgeTskErr := repository.RetrieveByTask(tsk.Task())
				if retPledgeTskErr == nil {
					str := fmt.Sprintf("the given pledge task (ID: %s) contains a Task (ID: %s) that is already assigned to another pledge Task", tsk.ID().String(), tsk.Task().ID().String())
					return errors.New(str)
				}

				// make sure the pledge is to the milestone wallet:
				if bytes.Compare(tsk.Pledge().To().ID().Bytes(), tsk.Task().Milestone().Wallet().ID().Bytes()) != 0 {
					str := fmt.Sprintf("the pledge (ID: %s) should be to the milestone (ID: %s) wallet (ID: %s), but it is to this wallet (ID: %s)", tsk.Pledge().ID().String(), tsk.Task().Milestone().ID().String(), tsk.Task().Milestone().Wallet().ID().String(), tsk.Pledge().To().ID().String())
					return errors.New(str)
				}

				// make sure the pledge is of the right amount:
				if tsk.Pledge().From().Amount() != tsk.Task().PledgeNeeded() {
					str := fmt.Sprintf("the task (ID: %s) requested %d tokens, %d pledged", tsk.Task().ID().String(), tsk.Task().PledgeNeeded(), tsk.Pledge().From().Amount())
					return errors.New(str)
				}

				// make sure the pledge does not exists:
				_, retPledgeErr := entityRepository.RetrieveByID(pledgeMetaData, tsk.Pledge().ID())
				if retPledgeErr == nil {
					str := fmt.Sprintf("the given pledge task (ID: %s) contains a Pledge (ID: %s) that already exists", tsk.ID().String(), tsk.Pledge().ID().String())
					return errors.New(str)
				}

				// save the pledge:
				savePledgeErr := entityService.Save(tsk.Pledge(), pledgeRepresentation)
				if savePledgeErr != nil {
					return savePledgeErr
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid pledge Task instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
