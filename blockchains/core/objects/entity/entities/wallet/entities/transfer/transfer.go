package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
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

func createTransferFromNormalized(normalized *normalizedTransfer) (Transfer, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	withIns, withInsErr := withdrawal.SDKFunc.CreateMetaData().Denormalize()(normalized.Withdrawal)
	if withInsErr != nil {
		return nil, withInsErr
	}

	depIns, depInsErr := deposit.SDKFunc.CreateMetaData().Denormalize()(normalized.Deposit)
	if depInsErr != nil {
		return nil, depInsErr
	}

	if with, ok := withIns.(withdrawal.Withdrawal); ok {
		if dep, ok := depIns.(deposit.Deposit); ok {
			out := createTransfer(&id, with, dep)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withIns.ID().String())
	return nil, errors.New(str)

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
	return obj.Dep
}
