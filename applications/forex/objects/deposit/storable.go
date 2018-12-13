package deposit

type storableDeposit struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	BankID string `json:"bank_id"`
}

func createStorableDeposit(ins Deposit) *storableDeposit {
	out := storableDeposit{
		ID:     ins.ID().String(),
		Amount: ins.Amount(),
		BankID: ins.Bank().ID().String(),
	}

	return &out
}
