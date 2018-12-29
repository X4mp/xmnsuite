package offer

type storableOffer struct {
	ID            string `json:"id"`
	PledgeID      string `json:"pledge_id"`
	ToID          string `json:"to_id"`
	Confirmations int    `json:"confirmations"`
	Amount        int    `json:"amount"`
	Price         int    `json:"price"`
	IP            string `json:"ip_address"`
	Port          int    `json:"port"`
}

func createStorableOffer(ins Offer) *storableOffer {
	out := storableOffer{
		ID:            ins.ID().String(),
		PledgeID:      ins.Pledge().ID().String(),
		ToID:          ins.To().ID().String(),
		Confirmations: ins.Confirmations(),
		Amount:        ins.Amount(),
		Price:         ins.Price(),
		IP:            ins.IP().String(),
		Port:          ins.Port(),
	}

	return &out
}
