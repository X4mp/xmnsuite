package external

type storableExternal struct {
	ID     string `json:"id"`
	LinkID string `json:"link_id"`
	ResID  string `json:"resource_id"`
}

func createStorableEWallet(ins External) *storableExternal {
	out := storableExternal{
		ID:     ins.ID().String(),
		LinkID: ins.Link().ID().String(),
		ResID:  ins.ResourceID().String(),
	}

	return &out
}
