package link

type normalizedLink struct {
	ID          string `json:"id"`
	Keyname     string `json:"keyname"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createNormalizedLink(link Link) (*normalizedLink, error) {
	out := normalizedLink{
		ID:          link.ID().String(),
		Title:       link.Title(),
		Description: link.Description(),
	}

	return &out, nil
}
