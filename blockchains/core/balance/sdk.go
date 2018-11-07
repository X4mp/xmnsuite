package balance

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

// Balance represents a wallet balance
type Balance interface {
	ID() *uuid.UUID
	On() wallet.Wallet
	Of() token.Token
	Amount() int
	CreatedOn() time.Time
}

// Repository represents a balance repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Balance, error)
	RetrieveByWalletAndToken(wal wallet.Wallet, tok token.Token) (Balance, error)
	RetrieveSetByWallet(wal wallet.Wallet) (entity.PartialSet, error)
	RetrieveSetByToken(tok token.Token) (entity.PartialSet, error)
}
