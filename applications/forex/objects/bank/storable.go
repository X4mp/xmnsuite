package bank

type storableBank struct {
	ID         string `json:"id"`
	PledgeID   string `json:"pledge_id"`
	CurrencyID string `json:"currency_id"`
	Amount     int    `json:"amount"`
	Price      int    `json:"price"`
}

func createStorableBank(ins Bank) *storableBank {
	out := storableBank{
		ID:         ins.ID().String(),
		PledgeID:   ins.Pledge().ID().String(),
		CurrencyID: ins.Currency().ID().String(),
		Amount:     ins.Amount(),
		Price:      ins.Price(),
	}

	return &out
}
