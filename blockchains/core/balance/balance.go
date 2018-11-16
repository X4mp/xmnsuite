package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

type balance struct {
	OnWallet wallet.Wallet `json:"on"`
	OfToken  token.Token   `json:"of"`
	Am       int           `json:"amount"`
}

func createBalance(on wallet.Wallet, of token.Token, amount int) Balance {
	out := balance{
		OnWallet: on,
		OfToken:  of,
		Am:       amount,
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
