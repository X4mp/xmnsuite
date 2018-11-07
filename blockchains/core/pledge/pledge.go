package pledge

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

type pledge struct {
	UUID           *uuid.UUID            `json:"id"`
	FromWithdrawal withdrawal.Withdrawal `json:"from"`
	ToWallet       wallet.Wallet         `json:"to"`
	Am             int                   `json:"amount"`
}

func createPledge(id *uuid.UUID, from withdrawal.Withdrawal, to wallet.Wallet, amount int) Pledge {
	out := pledge{
		UUID:           id,
		FromWithdrawal: from,
		ToWallet:       to,
		Am:             amount,
	}

	return &out
}

// ID returns the ID
func (obj *pledge) ID() *uuid.UUID {
	return obj.UUID
}

// From returns the from withdrawal
func (obj *pledge) From() withdrawal.Withdrawal {
	return obj.FromWithdrawal
}

// To returns the to wallet
func (obj *pledge) To() wallet.Wallet {
	return obj.ToWallet
}

// Amount returns the amount to pledge
func (obj *pledge) Amount() int {
	return obj.Am
}
