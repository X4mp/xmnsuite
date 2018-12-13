package fiatchain

import "github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"

type normalizedFiatChain struct {
	ID       string               `json:"id"`
	Seeds    []string             `json:"seeds"`
	Deposits []deposit.Normalized `json:"deposits"`
}

func createNormalizedFiatChain(ins FiatChain) (*normalizedFiatChain, error) {

	deps := ins.Deposits()
	normalizedDeps := []deposit.Normalized{}
	depNormalizeFunc := deposit.SDKFunc.CreateMetaData().Normalize()
	for _, oneDep := range deps {
		oneNormalized, oneNormalizedErr := depNormalizeFunc(oneDep)
		if oneNormalizedErr != nil {
			return nil, oneNormalizedErr
		}

		normalizedDeps = append(normalizedDeps, oneNormalized)
	}

	out := normalizedFiatChain{
		ID:       ins.ID().String(),
		Seeds:    ins.Seeds(),
		Deposits: normalizedDeps,
	}

	return &out, nil
}
