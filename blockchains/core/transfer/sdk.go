package transfer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	From() wallet.Wallet
	Token() token.Token
	Amount() int
	Content() string
	PubKey() crypto.PublicKey
}

// Service represents a transfer service
type Service interface {
	Save(ins Transfer) error
}
