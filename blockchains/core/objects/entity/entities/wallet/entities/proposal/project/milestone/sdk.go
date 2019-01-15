package milestone

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
)

// Milestone represents a project milestone
type Milestone interface {
	ID() *uuid.UUID
	Project() project.Project
	HasFeature() bool
	Feature() feature.Feature
	Wallet() wallet.Wallet
	Shares() int
	Title() string
	Description() string
	Details() string
}

// Normalized represents a normalized milestone
type Normalized interface {
}

// Repository represents a milestone repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Milestone, error)
	RetrieveByWallet(wal wallet.Wallet) (Milestone, error)
	RetrieveSetByProject(proj project.Project, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Project     project.Project
	Feature     feature.Feature
	Wallet      wallet.Wallet
	Shares      int
	Title       string
	Description string
	Details     string
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Milestone SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Milestone
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Milestone {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		if params.Feature != nil {
			out, outErr := createMilestoneWithFeature(params.ID, params.Project, params.Wallet, params.Shares, params.Title, params.Description, params.Details, params.Feature)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := createMilestone(params.ID, params.Project, params.Wallet, params.Shares, params.Title, params.Description, params.Details)
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
