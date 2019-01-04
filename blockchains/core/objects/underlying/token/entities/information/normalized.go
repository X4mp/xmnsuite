package information

type normalizedInformation struct {
	ID                   string `json:"id"`
	ConcensusNeeded      int    `json:"concensus_needed"`
	GzPricePerKb         int    `json:"gaz_price_per_kb"`
	MxAmountOfValidators int    `json:"max_amount_of_validators"`
}

func createNormalizedInformation(ins Information) (*normalizedInformation, error) {
	out := normalizedInformation{
		ID:                   ins.ID().String(),
		ConcensusNeeded:      ins.ConcensusNeeded(),
		GzPricePerKb:         ins.GazPricePerKb(),
		MxAmountOfValidators: ins.MaxAmountOfValidators(),
	}

	return &out, nil
}
