package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"

	"github.com/xmnservices/xmnsuite/datastore"
)

// Deposit represents the initial deposit
type Deposit interface {
	ID() *uuid.UUID
	To() wallet.Wallet
	Amount() int
}

// Normalized represents the normalized deposit
type Normalized interface {
}

// Repository represents the deposit Repository
type Repository interface {
	RetrieveSetByToWallet(wal wallet.Wallet) ([]Deposit, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	To     wallet.Wallet
	Amount int
}

// CreateRepositoryParams represents a CreateRepository params
type CreateRepositoryParams struct {
	Datastore        datastore.DataStore
	EntityRepository entity.Repository
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Deposit
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Deposit {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createDeposit(params.ID, params.To, params.Amount)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if deposit, ok := ins.(Deposit); ok {
					out := createStorableDeposit(deposit)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if deposit, ok := ins.(Deposit); ok {
					return []string{
						retrieveAllDepositsKeyname(),
						retrieveDepositsByToWalletIDKeyname(deposit.To().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)

			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {

				// metadata and representation:
				walletRepresentation := wallet.SDKFunc.CreateRepresentation()

				// create the entity repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				if deposit, ok := ins.(Deposit); ok {
					// try to retrieve the wallet:
					toWallet := deposit.To()
					_, retToWalletErr := repository.RetrieveByID(walletRepresentation.MetaData(), toWallet.ID())
					if retToWalletErr != nil {
						// save the wallet:
						saveErr := service.Save(toWallet, walletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		met := createMetaData()
		if params.Datastore != nil {
			entityRepository := entity.SDKFunc.CreateRepository(params.Datastore)
			out := createRepository(entityRepository, met)
			return out
		}

		out := createRepository(params.EntityRepository, met)
		return out
	},
}
