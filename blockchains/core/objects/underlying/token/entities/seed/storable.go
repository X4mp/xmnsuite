package seed

type storableSeed struct {
	ID     string `json:"id"`
	LinkID string `json:"link_id"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}

func createStorableSeed(ins Seed) *storableSeed {
	out := storableSeed{
		ID:     ins.ID().String(),
		LinkID: ins.Link().ID().String(),
		IP:     ins.IP().String(),
		Port:   ins.Port(),
	}

	return &out
}
