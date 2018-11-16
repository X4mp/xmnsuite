package pledge

type storablePledge struct {
	ID               string `json:"id"`
	FromWithdrawalID string `json:"from_withdrawal_id"`
	ToWalletID       string `json:"to_wallet_id"`
}

func createStorablePledge(ins Pledge) *storablePledge {
	out := storablePledge{
		ID:               ins.ID().String(),
		FromWithdrawalID: ins.From().ID().String(),
		ToWalletID:       ins.To().ID().String(),
	}

	return &out
}
