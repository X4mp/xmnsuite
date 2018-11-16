package pledge

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

type normalizedPledge struct {
	ID   string                `json:"id"`
	From withdrawal.Normalized `json:"from"`
	To   wallet.Normalized     `json:"to"`
}

func createNormalizedPledge(ins Pledge) (*normalizedPledge, error) {
	from, fromErr := withdrawal.SDKFunc.CreateMetaData().Normalize()(ins.From())
	if fromErr != nil {
		return nil, fromErr
	}

	to, toErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.To())
	if toErr != nil {
		return nil, toErr
	}

	out := normalizedPledge{
		ID:   ins.ID().String(),
		From: from,
		To:   to,
	}

	return &out, nil
}
