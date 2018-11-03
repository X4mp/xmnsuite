package withdrawal

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

type withdrawal struct {
	UUID       *uuid.UUID    `json:"id"`
	FromWallet wallet.Wallet `json:"from"`
	Am         int           `json:"amount"`
}

func createWithdrawal(id *uuid.UUID, fromWallet wallet.Wallet, amount int) Withdrawal {
	out := withdrawal{
		UUID:       id,
		FromWallet: fromWallet,
		Am:         amount,
	}

	return &out
}

// ID returns the ID
func (app *withdrawal) ID() *uuid.UUID {
	return app.UUID
}

// From returns the from wallet
func (app *withdrawal) From() wallet.Wallet {
	return app.FromWallet
}

// Amount returns the amount
func (app *withdrawal) Amount() int {
	return app.Am
}
