package affiliates

type storableAffiliate struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	URL     string `json:"url"`
}

func createStorableAffiliate(ins Affiliate) *storableAffiliate {
	out := storableAffiliate{
		ID:      ins.ID().String(),
		OwnerID: ins.Owner().ID().String(),
		URL:     ins.URL().String(),
	}

	return &out
}
