package chain

type storableChain struct {
	ID          string `json:"id"`
	OfferID     string `json:"offer_id"`
	TotalAmount int    `json:"total_amount"`
}

func createStorableChain(ins Chain) *storableChain {
	out := storableChain{
		ID:          ins.ID().String(),
		OfferID:     ins.Offer().ID().String(),
		TotalAmount: ins.TotalAmount(),
	}

	return &out
}
