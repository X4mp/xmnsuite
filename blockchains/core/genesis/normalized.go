package genesis

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
)

type normalizedGenesis struct {
	ID                   string             `json:"id"`
	GzPricePerKb         int                `json:"gaz_price_per_kb"`
	MxAmountOfValidators int                `json:"max_amount_of_validators"`
	Deposit              deposit.Normalized `json:"deposit"`
}

func createNormalizedGenesis(ins Genesis) (*normalizedGenesis, error) {
	normalizedDeposit, normalizedDepositErr := deposit.SDKFunc.CreateMetaData().Normalize()(ins.Deposit())
	if normalizedDepositErr != nil {
		return nil, normalizedDepositErr
	}

	out := normalizedGenesis{
		ID:                   ins.ID().String(),
		GzPricePerKb:         ins.GazPricePerKb(),
		MxAmountOfValidators: ins.MaxAmountOfValidators(),
		Deposit:              normalizedDeposit,
	}

	return &out, nil
}
