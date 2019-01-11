package project

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
)

// Project represents a project
type Project interface {
	ID() *uuid.UUID
	Project() approved_project.Project
	Owner() wallet.Wallet
	Manager() wallet.Wallet
	ManagerShares() int
	Linker() wallet.Wallet
	LinkerShares() int
	WorkerShares() int
}

// Normalized represents the normalized project
type Normalized interface {
}

// Repository represents the project repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Project, error)
	RetrieveByProject(proj approved_project.Project) (Project, error)
	RetrieveByOwner(owner wallet.Wallet) (Project, error)
	RetrieveByManager(mgr wallet.Wallet) (Project, error)
	RetrieveByLinker(linker wallet.Wallet) (Project, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID            *uuid.UUID
	Project       approved_project.Project
	Owner         wallet.Wallet
	Manager       wallet.Wallet
	ManagerShares int
	Linker        wallet.Wallet
	LinkerShares  int
	WorkerShares  int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Project SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Project
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Project {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createProject(
			params.ID,
			params.Project,
			params.Owner,
			params.Manager,
			params.ManagerShares,
			params.Linker,
			params.LinkerShares,
			params.WorkerShares,
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
