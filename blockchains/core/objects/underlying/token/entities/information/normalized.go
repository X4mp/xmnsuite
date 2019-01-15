package information

type normalizedInformation struct {
	ID                   string `json:"id"`
	ConcensusNeeded      int    `json:"concensus_needed"`
	GzPricePerKb         int    `json:"gaz_price_per_kb"`
	MxAmountOfValidators int    `json:"max_amount_of_validators"`
	NetworkShare         int    `json:"network_share"`
	ValidatorsShare      int    `json:"validator_share"`
	AffiliateShare       int    `json:"affiliate_share"`
}

func createNormalizedInformation(ins Information) (*normalizedInformation, error) {
	out := normalizedInformation{
		ID:                   ins.ID().String(),
		ConcensusNeeded:      ins.ConcensusNeeded(),
		GzPricePerKb:         ins.GazPricePerKb(),
		MxAmountOfValidators: ins.MaxAmountOfValidators(),
		NetworkShare:         ins.NetworkShare(),
		ValidatorsShare:      ins.ValidatorsShare(),
		AffiliateShare:       ins.AffiliateShare(),
	}

	return &out, nil
}
