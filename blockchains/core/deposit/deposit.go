package deposit

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

type deposit struct {
	UUID     *uuid.UUID    `json:"id"`
	ToWallet wallet.Wallet `json:"to"`
	Tok      token.Token   `json:"token"`
	Am       int           `json:"amount"`
}

func createDeposit(id *uuid.UUID, toWallet wallet.Wallet, tok token.Token, amount int) Deposit {
	out := deposit{
		UUID:     id,
		ToWallet: toWallet,
		Tok:      tok,
		Am:       amount,
	}

	return &out
}

// ID returns the ID
func (app *deposit) ID() *uuid.UUID {
	return app.UUID
}

// To returns the to user
func (app *deposit) To() wallet.Wallet {
	return app.ToWallet
}

// Token returns the token
func (app *deposit) Token() token.Token {
	return app.Tok
}

// Amount returns the amount
func (app *deposit) Amount() int {
	return app.Am
}
