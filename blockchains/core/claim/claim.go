package claim

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/transfer"
)

type claim struct {
	UUID *uuid.UUID        `json:"id"`
	Trx  transfer.Transfer `json:"transfer"`
	Dep  deposit.Deposit   `json:"deposit"`
}

func createClaim(id *uuid.UUID, trx transfer.Transfer, dep deposit.Deposit) Claim {
	out := claim{
		UUID: id,
		Trx:  trx,
		Dep:  dep,
	}

	return &out
}

// ID returns the ID
func (obj *claim) ID() *uuid.UUID {
	return obj.UUID
}

// Transfer returns the transfer
func (obj *claim) Transfer() transfer.Transfer {
	return obj.Trx
}

// Deposit returns the deposit
func (obj *claim) Deposit() deposit.Deposit {
	return obj.Dep
}
