package transfer

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

type transfer struct {
	UUID   *uuid.UUID            `json:"id"`
	Withdr withdrawal.Withdrawal `json:"withdrawal"`
	Dep    deposit.Deposit       `json:"deposit"`
}

func createTransfer(id *uuid.UUID, withdrawal withdrawal.Withdrawal, dep deposit.Deposit) Transfer {
	out := transfer{
		UUID:   id,
		Withdr: withdrawal,
		Dep:    dep,
	}

	return &out
}

// ID returns the ID
func (obj *transfer) ID() *uuid.UUID {
	return obj.UUID
}

// Withdrawal returns the from withdrawal
func (obj *transfer) Withdrawal() withdrawal.Withdrawal {
	return obj.Withdr
}

// Deposit returns the deposit
func (obj *transfer) Deposit() deposit.Deposit {
	return nil
}
