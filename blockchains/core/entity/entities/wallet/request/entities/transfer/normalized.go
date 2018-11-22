package transfer

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

type normalizedTransfer struct {
	ID         string                `json:"id"`
	Withdrawal withdrawal.Normalized `json:"withdrawal"`
	Deposit    deposit.Normalized    `json:"deposit"`
}

func createNormalizedTransfer(trx Transfer) (*normalizedTransfer, error) {

	with, withErr := withdrawal.SDKFunc.CreateMetaData().Normalize()(trx.Withdrawal())
	if withErr != nil {
		return nil, withErr
	}

	dep, depErr := deposit.SDKFunc.CreateMetaData().Normalize()(trx.Deposit())
	if depErr != nil {
		return nil, depErr
	}

	out := normalizedTransfer{
		ID:         trx.ID().String(),
		Withdrawal: with,
		Deposit:    dep,
	}

	return &out, nil
}
