package developer

type storableDeveloper struct {
	ID     string `json:"id"`
	UserID string `json:"user"`
	Name   string `json:"name"`
	Resume string `json:"resume"`
}

func createStorableDeveloper(dev Developer) *storableDeveloper {
	out := storableDeveloper{
		ID:     dev.ID().String(),
		UserID: dev.User().ID().String(),
		Name:   dev.Name(),
		Resume: dev.Resume(),
	}

	return &out
}
