package chain

import (
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
)

type normalizedChain struct {
	ID          string           `json:"id"`
	Offer       offer.Normalized `json:"offer"`
	TotalAmount int              `json:"total_amount"`
}

func createNormalizedChain(ins Chain) (*normalizedChain, error) {
	off, offErr := offer.SDKFunc.CreateMetaData().Normalize()(ins.Offer())
	if offErr != nil {
		return nil, offErr
	}

	out := normalizedChain{
		ID:          ins.ID().String(),
		Offer:       off,
		TotalAmount: ins.TotalAmount(),
	}

	return &out, nil

}
