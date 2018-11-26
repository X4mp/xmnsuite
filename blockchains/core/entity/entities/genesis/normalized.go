package genesis

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
)

type normalizedGenesis struct {
	ID                   string             `json:"id"`
	ConcensusNeeded      int                `json:"concensus_needed"`
	GzPricePerKb         int                `json:"gaz_price_per_kb"`
	MxAmountOfValidators int                `json:"max_amount_of_validators"`
	User                 user.Normalized    `json:"user"`
	Deposit              deposit.Normalized `json:"deposit"`
}

func createNormalizedGenesis(ins Genesis) (*normalizedGenesis, error) {
	normalizedDeposit, normalizedDepositErr := deposit.SDKFunc.CreateMetaData().Normalize()(ins.Deposit())
	if normalizedDepositErr != nil {
		return nil, normalizedDepositErr
	}

	normalizedUser, normalizedUserErr := user.SDKFunc.CreateMetaData().Normalize()(ins.User())
	if normalizedUserErr != nil {
		return nil, normalizedUserErr
	}

	out := normalizedGenesis{
		ID:                   ins.ID().String(),
		ConcensusNeeded:      ins.ConcensusNeeded(),
		GzPricePerKb:         ins.GazPricePerKb(),
		MxAmountOfValidators: ins.MaxAmountOfValidators(),
		User:                 normalizedUser,
		Deposit:              normalizedDeposit,
	}

	return &out, nil
}
