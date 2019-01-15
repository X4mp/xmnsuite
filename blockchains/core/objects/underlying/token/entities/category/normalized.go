package category

type normalizedCategory struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Parent      Normalized `json:"parent_category"`
}

func createNormalizedCategory(ins Category) (*normalizedCategory, error) {
	var parent Normalized
	if ins.HasParent() {
		par, parErr := createMetaData().Normalize()(ins.Parent())
		if parErr != nil {
			return nil, parErr
		}

		parent = par
	}

	out := normalizedCategory{
		ID:          ins.ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
		Parent:      parent,
	}

	return &out, nil
}
