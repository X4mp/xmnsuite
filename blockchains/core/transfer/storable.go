package transfer

type storableTransfer struct {
	ID           string `json:"id"`
	WithdrawalID string `json:"withdrawal_id"`
	Content      string `json:"content"`
	PubKey       string `json:"public_key"`
}

func createStorableTransfer(trans Transfer) *storableTransfer {
	out := storableTransfer{
		ID:           trans.ID().String(),
		WithdrawalID: trans.Withdrawal().ID().String(),
		Content:      trans.Content(),
		PubKey:       trans.PubKey().String(),
	}

	return &out
}
