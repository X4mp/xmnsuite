package genesis

type storableGenesis struct {
	ID                       string `json:"id"`
	ConcensusNeeded          int    `json:"concensus_needed"`
	DeveloperConcensusNeeded int    `json:"developer_concensus_needed"`
	GzPricePerKb             int    `json:"gaz_price_per_kb"`
	MxAmountOfValidators     int    `json:"max_amount_of_validators"`
	UserID                   string `json:"user_id"`
	InitialDepositID         string `json:"initial_deposit_id"`
}

func createStorableGenesis(gen Genesis) *storableGenesis {
	out := storableGenesis{
		ID:                       gen.ID().String(),
		ConcensusNeeded:          gen.ConcensusNeeded(),
		DeveloperConcensusNeeded: gen.DeveloperConcensusNeeded(),
		GzPricePerKb:             gen.GazPricePerKb(),
		MxAmountOfValidators:     gen.MaxAmountOfValidators(),
		UserID:                   gen.User().ID().String(),
		InitialDepositID:         gen.Deposit().ID().String(),
	}

	return &out
}
