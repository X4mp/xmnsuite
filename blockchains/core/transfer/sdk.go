package transfer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Withdrawal() withdrawal.Withdrawal
	Content() string
	PubKey() crypto.PublicKey
}

// Service represents a transfer service
type Service interface {
	Save(ins Transfer) error
}
