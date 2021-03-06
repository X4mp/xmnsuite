package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Withdrawal represents a withdrawal
type Withdrawal interface {
	ID() *uuid.UUID
	From() wallet.Wallet
	Amount() int
}

// Normalized represents the normalized withdrawal
type Normalized interface {
}

// Repository represents the withdrawal repository
type Repository interface {
	RetrieveSetByFromWallet(wal wallet.Wallet) ([]Withdrawal, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	From   wallet.Wallet
	Amount int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	Datastore        datastore.DataStore
	EntityRepository entity.Repository
}

// SDKFunc represents the Withdrawal SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Withdrawal
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Withdrawal {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createWithdrawal(params.ID, params.From, params.Amount)
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
				if withdrawal, ok := ins.(Withdrawal); ok {
					out := createStorableWithdrawal(withdrawal)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if with, ok := ins.(Withdrawal); ok {
					return []string{
						retrieveAllWithdrawalsKeyname(),
						retrieveWithdrawalsByFromWalletIDKeyname(with.From().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)

			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {

				calculate := func(withs []Withdrawal, deps []deposit.Deposit) (int, error) {
					// calculate the withdrawals amount:
					withAmount := 0
					for _, oneWithdrawalIns := range withs {
						withAmount += oneWithdrawalIns.(Withdrawal).Amount()
					}

					// calculate the deposits amount:
					depAmount := 0
					for _, oneDepIns := range deps {
						depAmount += oneDepIns.(deposit.Deposit).Amount()
					}

					// create the balance:
					total := depAmount - withAmount
					return total, nil
				}

				if with, ok := ins.(Withdrawal); ok {

					//create the metadata and representation:
					walletRepresentation := wallet.SDKFunc.CreateRepresentation()
					metaData := createMetaData()

					// create the repositories:
					repository := entity.SDKFunc.CreateRepository(ds)
					depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
						Datastore: ds,
					})
					withdrawalRepository := createRepository(repository, metaData)

					// make sure the withdrawal does not already exists:
					_, withErr := repository.RetrieveByID(metaData, with.ID())
					if withErr == nil {
						str := fmt.Sprintf("the Withdrawal instance (ID: %s) already exists", with.ID().String())
						return errors.New(str)
					}

					// fetch the wallet:
					wal := with.From()

					// try to retrieve the wallet:
					_, retToWalletErr := repository.RetrieveByID(walletRepresentation.MetaData(), wal.ID())
					if retToWalletErr != nil {
						str := fmt.Sprintf("the Withdrawal instance (ID: %s) contains a Wallet instance (ID: %s) that do not exists", with.ID().String(), wal.ID().String())
						return errors.New(str)
					}

					// retrieve all the withdrawals related to our wallet:
					withs, withsErr := withdrawalRepository.RetrieveSetByFromWallet(wal)
					if withsErr != nil {
						return withsErr
					}

					// retrieve all the deposits related to our wallet:
					deps, depsErr := depositRepository.RetrieveSetByToWallet(wal)
					if depsErr != nil {
						return depsErr
					}

					// retrieve the balance:
					balance, balanceErr := calculate(withs, deps)
					if balanceErr != nil {
						str := fmt.Sprintf("there was an error while retrieving the balance of the Wallet (ID: %s): %s", with.From().ID().String(), balanceErr.Error())
						return errors.New(str)
					}

					// make sure the balance is bigger or equal to the withdrawal:
					if balance <= with.Amount() {
						str := fmt.Sprintf("the balance of the wallet (ID: %s) is %d, but the transfer needed %d", with.From().ID().String(), balance, with.Amount())
						return errors.New(str)
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
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
