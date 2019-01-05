package transfer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Withdrawal() withdrawal.Withdrawal
	Deposit() deposit.Deposit
}

// Normalized represents the normalized transfer
type Normalized interface {
}

// Repository represents the transfer reposiotry
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Transfer, error)
	RetrieveByDeposit(dep deposit.Deposit) (Transfer, error)
	RetrieveByWithdrawal(with withdrawal.Withdrawal) (Transfer, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Withdrawal withdrawal.Withdrawal
	Deposit    deposit.Deposit
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Transfer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Transfer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Transfer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createTransfer(params.ID, params.Withdrawal, params.Deposit)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
}
