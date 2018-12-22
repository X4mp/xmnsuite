package group

type storableGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createStorableGroup(ins Group) *storableGroup {
	out := storableGroup{
		ID:   ins.ID().String(),
		Name: ins.Name(),
	}

	return &out
}
