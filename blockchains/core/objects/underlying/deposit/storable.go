package deposit

type storableDeposit struct {
	ID         string `json:"id"`
	ToWalletID string `json:"to_wallet_id"`
	TokenID    string `json:"token_id"`
	Amount     int    `json:"amount"`
}

func createStorableDeposit(ins Deposit) *storableDeposit {
	out := storableDeposit{
		ID:         ins.ID().String(),
		ToWalletID: ins.To().ID().String(),
		TokenID:    ins.Token().ID().String(),
		Amount:     ins.Amount(),
	}

	return &out
}
