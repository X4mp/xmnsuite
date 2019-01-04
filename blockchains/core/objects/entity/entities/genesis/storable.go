package genesis

type storableGenesis struct {
	ID               string `json:"id"`
	InfoID           string `json:"information_id"`
	UserID           string `json:"user_id"`
	InitialDepositID string `json:"initial_deposit_id"`
}

func createStorableGenesis(gen Genesis) *storableGenesis {
	out := storableGenesis{
		ID:               gen.ID().String(),
		InfoID:           gen.Info().ID().String(),
		UserID:           gen.User().ID().String(),
		InitialDepositID: gen.Deposit().ID().String(),
	}

	return &out
}
