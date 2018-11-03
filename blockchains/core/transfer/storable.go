package transfer

type storableTransfer struct {
	ID           string `json:"id"`
	FromWalletID string `json:"from_wallet_id"`
	TokenID      string `json:"token_id"`
	Amount       int    `json:"amount"`
	Content      string `json:"content"`
	PubKey       string `json:"public_key"`
}

func createStorableTransfer(trans Transfer) *storableTransfer {
	out := storableTransfer{
		ID:           trans.ID().String(),
		FromWalletID: trans.From().ID().String(),
		TokenID:      trans.Token().ID().String(),
		Amount:       trans.Amount(),
		Content:      trans.Content(),
		PubKey:       trans.PubKey().String(),
	}

	return &out
}
