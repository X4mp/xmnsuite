package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
)

type pledge struct {
	UUID           *uuid.UUID            `json:"id"`
	FromWithdrawal withdrawal.Withdrawal `json:"from"`
	ToWallet       wallet.Wallet         `json:"to"`
}

func createPledge(id *uuid.UUID, from withdrawal.Withdrawal, to wallet.Wallet) Pledge {
	out := pledge{
		UUID:           id,
		FromWithdrawal: from,
		ToWallet:       to,
	}

	return &out
}

func createPledgeFromNormalized(normalized *normalizedPledge) (Pledge, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	fromIns, fromInsErr := withdrawal.SDKFunc.CreateMetaData().Denormalize()(normalized.From)
	if fromInsErr != nil {
		return nil, fromInsErr
	}

	toIns, toInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.To)
	if toInsErr != nil {
		return nil, toInsErr
	}

	if from, ok := fromIns.(withdrawal.Withdrawal); ok {
		if to, ok := toIns.(wallet.Wallet); ok {
			out := createPledge(&id, from, to)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", toIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", fromIns.ID().String())
	return nil, errors.New(str)

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
