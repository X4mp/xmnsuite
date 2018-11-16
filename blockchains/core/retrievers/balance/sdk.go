package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Balance represents a wallet balance
type Balance interface {
	On() wallet.Wallet
	Of() token.Token
	Amount() int
}

// Repository represents a balance repository
type Repository interface {
	RetrieveByWalletAndToken(wal wallet.Wallet, tok token.Token) (Balance, error)
}

// SDKFunc represents the balance SDK func
var SDKFunc = struct {
	CreateRepository func(ds datastore.DataStore) Repository
}{
	CreateRepository: func(ds datastore.DataStore) Repository {
		depositRepository := deposit.SDKFunc.CreateRepository(ds)
		withdrawalRepository := withdrawal.SDKFunc.CreateRepository(ds)
		out := createRepository(depositRepository, withdrawalRepository)
		return out
	},
}