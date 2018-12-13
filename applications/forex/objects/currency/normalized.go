package currency

import (
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
)

type normalizedCurrency struct {
	ID          string              `json:"id"`
	Category    category.Normalized `json:"category"`
	Symbol      string              `json:"symbol"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
}

func createNormalizedCurrency(ins Currency) (*normalizedCurrency, error) {
	cat, catErr := category.SDKFunc.CreateMetaData().Normalize()(ins.Category())
	if catErr != nil {
		return nil, catErr
	}

	out := normalizedCurrency{
		ID:          ins.ID().String(),
		Category:    cat,
		Symbol:      ins.Symbol(),
		Name:        ins.Name(),
		Description: ins.Description(),
	}

	return &out, nil
}
