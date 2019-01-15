package task

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

// Task represents a task
type Task interface {
	ID() *uuid.UUID
	Milestone() milestone.Milestone
	CreatedBy() user.User
	Title() string
	Details() string
	Deadline() time.Time
	Reward() int
	PledgeNeeded() int
}

// Normalized represents a normalized task
type Normalized interface {
}

// Repository represents a Task repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Task, error)
	RetrieveSetByMilestone(mils milestone.Milestone, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByCreatedByUser(crBy user.User, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID           *uuid.UUID
	Milestone    milestone.Milestone
	CreatedBy    user.User
	Title        string
	Details      string
	Deadline     time.Time
	Reward       int
	PledgeNeeded int
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

		out, outErr := createTask(
			params.ID,
			params.Milestone,
			params.CreatedBy,
			params.Title,
			params.Details,
			params.Deadline,
			params.Reward,
			params.PledgeNeeded,
		)

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
