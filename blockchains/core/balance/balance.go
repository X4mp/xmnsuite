package balance

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

type balance struct {
	UUID     *uuid.UUID    `json:"id"`
	OnWallet wallet.Wallet `json:"on"`
	OfToken  token.Token   `json:"of"`
	Am       int           `json:"amount"`
	CrOn     time.Time     `json:"created_on"`
}

func createBalance(id *uuid.UUID, on wallet.Wallet, of token.Token, amount int, createdOn time.Time) Balance {
	out := balance{
		UUID:     id,
		OnWallet: on,
		OfToken:  of,
		Am:       amount,
		CrOn:     createdOn,
	}

	return &out
}

// ID returns the ID
func (obj *balance) ID() *uuid.UUID {
	return obj.UUID
}

// On returns the on wallet
func (obj *balance) On() wallet.Wallet {
	return obj.OnWallet
}

// Of returns the of token
func (obj *balance) Of() token.Token {
	return obj.OfToken
}

// Amount returns the amount
func (obj *balance) Amount() int {
	return obj.Am
}

// CreatedOn returns the creation time
func (obj *balance) CreatedOn() time.Time {
	return obj.CrOn
}
