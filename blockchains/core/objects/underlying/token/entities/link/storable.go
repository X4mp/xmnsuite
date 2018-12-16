package link

type storableLink struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createStorableLink(ins Link) *storableLink {
	out := storableLink{
		ID:          ins.ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
	}

	return &out
}
