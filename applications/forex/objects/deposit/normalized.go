package deposit

import (
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
)

type normalizedDeposit struct {
	ID     string          `json:"id"`
	Amount int             `json:"amount"`
	Bank   bank.Normalized `json:"bank"`
}

func createNormalizedDeposit(ins Deposit) (*normalizedDeposit, error) {
	bnk, bnkErr := bank.SDKFunc.CreateMetaData().Normalize()(ins.Bank())
	if bnkErr != nil {
		return nil, bnkErr
	}

	out := normalizedDeposit{
		ID:     ins.ID().String(),
		Amount: ins.Amount(),
		Bank:   bnk,
	}

	return &out, nil
}
