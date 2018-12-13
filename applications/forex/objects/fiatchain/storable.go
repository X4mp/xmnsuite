package fiatchain

type storableFiatChain struct {
	ID         string   `json:"id"`
	Seeds      []string `json:"seeds"`
	DepositIDs []string `json:"deposit_ids"`
}

func createStorableFiatChain(ins FiatChain) *storableFiatChain {
	depIDs := []string{}
	deps := ins.Deposits()
	for _, oneDep := range deps {
		depIDs = append(depIDs, oneDep.ID().String())
	}

	out := storableFiatChain{
		ID:         ins.ID().String(),
		Seeds:      ins.Seeds(),
		DepositIDs: depIDs,
	}

	return &out
}
