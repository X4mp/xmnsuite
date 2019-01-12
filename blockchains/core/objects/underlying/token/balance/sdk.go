package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Balance represents a wallet balance
type Balance interface {
	Wallet() wallet.Wallet
	Amount() int
}

// Repository represents a balance repository
type Repository interface {
	RetrieveByWallet(wal wallet.Wallet) (Balance, error)
}

// CreateRepositoryParams represents a CreateRepository params
type CreateRepositoryParams struct {
	Datastore            datastore.DataStore
	DepositRepository    deposit.Repository
	WithdrawalRepository withdrawal.Repository
}

// SDKFunc represents the balance SDK func
var SDKFunc = struct {
	CreateRepository func(params CreateRepositoryParams) Repository
}{
	CreateRepository: func(params CreateRepositoryParams) Repository {
		if params.Datastore != nil {
			if params.Datastore != nil {
				depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
					Datastore: params.Datastore,
				})

				withdrawalRepository := withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
					Datastore: params.Datastore,
				})

				out := createRepository(depositRepository, withdrawalRepository)
				return out
			}

			out := createRepository(params.DepositRepository, params.WithdrawalRepository)
			return out
		}

		out := createRepository(params.DepositRepository, params.WithdrawalRepository)
		return out
	},
}
