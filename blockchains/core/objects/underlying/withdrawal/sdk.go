package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Withdrawal represents a withdrawal
type Withdrawal interface {
	ID() *uuid.UUID
	From() wallet.Wallet
	Token() token.Token
	Amount() int
}

// Normalized represents the normalized withdrawal
type Normalized interface {
}

// Repository represents the withdrawal repository
type Repository interface {
	RetrieveSetByFromWalletAndToken(wal wallet.Wallet, tok token.Token) (entity.PartialSet, error)
}

// Data represents human-redable data
type Data struct {
	ID     string
	From   *wallet.Data
	Token  *token.Data
	Amount int
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Withdrawals []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	From   wallet.Wallet
	Token  token.Token
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
	ToData               func(with Withdrawal) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Withdrawal {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createWithdrawal(params.ID, params.From, params.Token, params.Amount)
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
				if withdrawal, ok := ins.(Withdrawal); ok {
					base := retrieveAllWithdrawalsKeyname()
					return []string{
						base,
						retrieveWithdrawalsByTokenIDKeyname(withdrawal.Token().ID()),
						retrieveWithdrawalsByFromWalletIDKeyname(withdrawal.From().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {

				calculate := func(withsPS entity.PartialSet, depsPS entity.PartialSet) (int, error) {
					// calculate the withdrawals amount:
					withAmount := 0
					withs := withsPS.Instances()
					for _, oneWithdrawalIns := range withs {
						withAmount += oneWithdrawalIns.(Withdrawal).Amount()
					}

					// calculate the deposits amount:
					depAmount := 0
					deps := depsPS.Instances()
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

					// fetch the wallet and token:
					wal := with.From()
					tok := with.Token()

					// try to retrieve the wallet:
					_, retToWalletErr := repository.RetrieveByID(walletRepresentation.MetaData(), wal.ID())
					if retToWalletErr != nil {
						str := fmt.Sprintf("the Withdrawal instance (ID: %s) contains a Wallet instance (ID: %s) that do not exists", with.ID().String(), wal.ID().String())
						return errors.New(str)
					}

					// retrieve all the withdrawals related to our wallet and token:
					withsPS, withsPSErr := withdrawalRepository.RetrieveSetByFromWalletAndToken(wal, tok)
					if withsPSErr != nil {
						return withsPSErr
					}

					// retrieve all the deposits related to our wallet and token:
					depsPS, depsPSErr := depositRepository.RetrieveSetByToWalletAndToken(wal, tok)
					if depsPSErr != nil {
						return depsPSErr
					}

					// retrieve the balance:
					balance, balanceErr := calculate(withsPS, depsPS)
					if balanceErr != nil {
						str := fmt.Sprintf("there was an error while retrieving the balance of the Wallet (ID: %s), for the Token (ID: %s): %s", with.From().ID().String(), with.Token().ID().String(), balanceErr.Error())
						return errors.New(str)
					}

					// make sure the balance is bigger or equal to the withdrawal:
					if balance <= with.Amount() {
						str := fmt.Sprintf("the balance of the wallet (ID: %s) for the token (ID: %s) is %d, but the transfer needed %d", with.From().ID().String(), with.Token().ID().String(), balance, with.Amount())
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
	ToData: func(with Withdrawal) *Data {
		return toData(with)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
