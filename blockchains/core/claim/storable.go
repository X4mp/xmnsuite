package claim

type storableClaim struct {
	ID         string `json:"id"`
	TransferID string `json:"transfer_id"`
	DepositID  string `json:"deposit_id"`
}

func createStorableClaim(ins Claim) *storableClaim {
	out := storableClaim{
		ID:         ins.ID().String(),
		TransferID: ins.Transfer().ID().String(),
		DepositID:  ins.Deposit().ID().String(),
	}

	return &out
}
