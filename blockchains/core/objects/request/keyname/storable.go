package keyname

type storableKeyname struct {
	ID      string `json:"id"`
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
}

func createStorableKeyname(ins Keyname) *storableKeyname {
	out := storableKeyname{
		ID:      ins.ID().String(),
		GroupID: ins.Group().ID().String(),
		Name:    ins.Name(),
	}

	return &out
}
