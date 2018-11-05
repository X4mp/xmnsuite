package withdrawal

type storableWithdrawal struct {
	ID           string `json:"id"`
	FromWalletID string `json:"from_wallet_id"`
	TokenID      string `json:"token_id"`
	Amount       int    `json:"amount"`
}

func createStorableWithdrawal(ins Withdrawal) *storableWithdrawal {
	out := storableWithdrawal{
		ID:           ins.ID().String(),
		FromWalletID: ins.From().ID().String(),
		TokenID:      ins.Token().ID().String(),
		Amount:       ins.Amount(),
	}

	return &out
}
