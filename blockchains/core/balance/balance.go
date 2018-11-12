package balance

import (
	"time"

	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

type balance struct {
	OnWallet wallet.Wallet `json:"on"`
	OfToken  token.Token   `json:"of"`
	Am       int           `json:"amount"`
	CrOn     time.Time     `json:"created_on"`
}

func createBalance(on wallet.Wallet, of token.Token, amount int, createdOn time.Time) Balance {
	out := balance{
		OnWallet: on,
		OfToken:  of,
		Am:       amount,
		CrOn:     createdOn,
	}

	return &out
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
