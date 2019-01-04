package information

type storableInformation struct {
	ID                   string `json:"id"`
	ConcensusNeeded      int    `json:"concensus_needed"`
	GzPricePerKb         int    `json:"gaz_price_per_kb"`
	MxAmountOfValidators int    `json:"max_amount_of_validators"`
}

func createStorableInformation(info Information) *storableInformation {
	out := storableInformation{
		ID:                   info.ID().String(),
		ConcensusNeeded:      info.ConcensusNeeded(),
		GzPricePerKb:         info.GazPricePerKb(),
		MxAmountOfValidators: info.MaxAmountOfValidators(),
	}

	return &out
}
