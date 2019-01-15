package withdrawal

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedWithdrawal struct {
	ID     string            `json:"id"`
	From   wallet.Normalized `json:"from"`
	Amount int               `json:"amount"`
}

func createNormalizedWithdrawal(ins Withdrawal) (*normalizedWithdrawal, error) {
	from, fromErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.From())
	if fromErr != nil {
		return nil, fromErr
	}

	out := normalizedWithdrawal{
		ID:     ins.ID().String(),
		From:   from,
		Amount: ins.Amount(),
	}

	return &out, nil
}
