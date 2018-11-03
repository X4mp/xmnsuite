package deposit

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

type deposit struct {
	UUID     *uuid.UUID    `json:"id"`
	ToWallet wallet.Wallet `json:"to"`
	Am       int           `json:"amount"`
}

func createDeposit(id *uuid.UUID, toWallet wallet.Wallet, amount int) Deposit {
	out := deposit{
		UUID:     id,
		ToWallet: toWallet,
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

// Amount returns the amount
func (app *deposit) Amount() int {
	return app.Am
}
