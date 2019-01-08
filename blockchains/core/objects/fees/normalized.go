package fees

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

type normalizedFee struct {
	ID        string                `json:"id"`
	Client    withdrawal.Normalized `json:"client"`
	Network   deposit.Normalized    `json:"network"`
	Validators     []deposit.Normalized  `json:"validators"`
	Affiliate deposit.Normalized    `json:"affiliate"`
}

func createNormalizedFee(ins Fee) (*normalizedFee, error) {
	client, clientErr := withdrawal.SDKFunc.CreateMetaData().Normalize()(ins.Client())
	if clientErr != nil {
		return nil, clientErr
	}

	normDepFunc := deposit.SDKFunc.CreateMetaData().Normalize()
	net, netErr := normDepFunc(ins.Network())
	if netErr != nil {
		return nil, netErr
	}

	validators := ins.Validators()
	normalizedValidators := []deposit.Normalized{}
	for _, oneValidator := range validators {
		nde, ndeErr := normDepFunc(oneValidator)
		if ndeErr != nil {
			return nil, ndeErr
		}

		normalizedValidators = append(normalizedValidators, nde)
	}

	var aff deposit.Normalized
	if ins.HasAffiliate() {
		normAff, normAffErr := normDepFunc(ins.Affiliate())
		if normAffErr != nil {
			return nil, normAffErr
		}

		aff = normAff
	}

	out := normalizedFee{
		ID:        ins.ID().String(),
		Client:    client,
		Network:   net,
		Validators:     normalizedValidators,
		Affiliate: aff,
	}

	return &out, nil
}
