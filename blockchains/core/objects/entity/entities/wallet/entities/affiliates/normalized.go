package affiliates

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedAffiliate struct {
	ID    string            `json:"id"`
	Owner wallet.Normalized `json:"owner"`
	URL   string            `json:"url"`
}

func createNormalizedAffiliate(ins Affiliate) (*normalizedAffiliate, error) {
	own, ownErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Owner())
	if ownErr != nil {
		return nil, ownErr
	}

	out := normalizedAffiliate{
		ID:    ins.ID().String(),
		Owner: own,
		URL:   ins.URL().String(),
	}

	return &out, nil
}
