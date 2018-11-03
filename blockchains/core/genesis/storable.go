package genesis

type storableGenesis struct {
	ID                   string `json:"id"`
	GzPricePerKb         int    `json:"gaz_price_per_kb"`
	MxAmountOfValidators int    `json:"max_amount_of_validators"`
	InitialDepositID     string `json:"initial_deposit_id"`
	TokenID              string `json:"token_id"`
}

func createStorableGenesis(gen Genesis) *storableGenesis {
	out := storableGenesis{
		ID:                   gen.ID().String(),
		GzPricePerKb:         gen.GazPricePerKb(),
		MxAmountOfValidators: gen.MaxAmountOfValidators(),
		InitialDepositID:     gen.Deposit().ID().String(),
		TokenID:              gen.Token().ID().String(),
	}

	return &out
}
