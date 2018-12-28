package deposit

type storableDeposit struct {
	ID      string `json:"id"`
	OfferID string `json:"offer_id"`
	FromID  string `json:"from_id"`
	Amount  int    `json:"amount"`
}

func createStorableDeposit(ins Deposit) *storableDeposit {
	out := storableDeposit{
		ID:      ins.ID().String(),
		OfferID: ins.Offer().ID().String(),
		FromID:  ins.From().ID().String(),
		Amount:  ins.Amount(),
	}

	return &out
}
