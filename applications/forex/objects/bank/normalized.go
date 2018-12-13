package bank

import (
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
)

type normalizedBank struct {
	ID       string              `json:"id"`
	Pledge   pledge.Normalized   `json:"pledge"`
	Currency currency.Normalized `json:"currency"`
	Amount   int                 `json:"amount"`
	Price    int                 `json:"price"`
}

func createNormalizedBank(ins Bank) (*normalizedBank, error) {
	pldge, pldgeErr := pledge.SDKFunc.CreateMetaData().Normalize()(ins.Pledge())
	if pldgeErr != nil {
		return nil, pldgeErr
	}

	curr, currErr := currency.SDKFunc.CreateMetaData().Normalize()(ins.Currency())
	if currErr != nil {
		return nil, currErr
	}

	out := normalizedBank{
		ID:       ins.ID().String(),
		Pledge:   pldge,
		Currency: curr,
		Amount:   ins.Amount(),
		Price:    ins.Price(),
	}

	return &out, nil
}
