package deposit

import (
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
)

type normalizedDeposit struct {
	ID     string             `json:"id"`
	Offer  offer.Normalized   `json:"offer"`
	From   address.Normalized `json:"from"`
	Amount int                `json:"amount"`
}

func createNormalizedDeposit(ins Deposit) (*normalizedDeposit, error) {

	off, offErr := offer.SDKFunc.CreateMetaData().Normalize()(ins.Offer())
	if offErr != nil {
		return nil, offErr
	}

	frm, frmErr := address.SDKFunc.CreateMetaData().Normalize()(ins.From())
	if frmErr != nil {
		return nil, frmErr
	}

	out := normalizedDeposit{
		ID:     ins.ID().String(),
		Offer:  off,
		From:   frm,
		Amount: ins.Amount(),
	}

	return &out, nil
}
