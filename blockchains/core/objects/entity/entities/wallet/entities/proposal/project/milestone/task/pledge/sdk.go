package pledge

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

// Task represents a pledge task
type Task interface {
	ID() *uuid.UUID
	Task() mils_task.Task
	Pledge() pledge.Pledge
}

// Normalized represents a normalized task
type Normalized interface {
}

// Repository represents a pledge task repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Task, error)
	RetrieveByTask(tsk mils_task.Task) (Task, error)
	RetrieveByPledge(pldge pledge.Pledge) (Task, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	Task   mils_task.Task
	Pledge pledge.Pledge
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Milestone SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Task
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Task {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createTask(params.ID, params.Task, params.Pledge)

		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return representation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
}
