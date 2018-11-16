package transfer

type storableTransfer struct {
	ID           string `json:"id"`
	WithdrawalID string `json:"withdrawal_id"`
	DepositID    string `json:"deposit_id"`
}

func createStorableTransfer(trans Transfer) *storableTransfer {
	out := storableTransfer{
		ID:           trans.ID().String(),
		WithdrawalID: trans.Withdrawal().ID().String(),
		DepositID:    trans.Deposit().ID().String(),
	}

	return &out
}
