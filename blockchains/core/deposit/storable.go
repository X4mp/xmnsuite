package deposit

type storableDeposit struct {
	ID         string `json:"id"`
	ToWalletID string `json:"to_wallet_id"`
	Amount     int    `json:"amount"`
}

func createStorableDeposit(ins Deposit) *storableDeposit {
	out := storableDeposit{
		ID:         ins.ID().String(),
		ToWalletID: ins.To().ID().String(),
		Amount:     ins.Amount(),
	}

	return &out
}
