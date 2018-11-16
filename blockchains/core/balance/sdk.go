package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
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
