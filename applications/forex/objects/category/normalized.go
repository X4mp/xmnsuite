package category

type normalizedCategory struct {
	ID          string              `json:"id"`
	Parent      *normalizedCategory `json:"parent"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
}

func createNormalizedCategory(ins Category) (*normalizedCategory, error) {
	var normalizedPar *normalizedCategory
	if ins.HasParent() {
		par, parErr := createNormalizedCategory(ins.Parent())
		if parErr != nil {
			return nil, parErr
		}

		normalizedPar = par
	}

	out := normalizedCategory{
		ID:          ins.ID().String(),
		Parent:      normalizedPar,
		Name:        ins.Name(),
		Description: ins.Description(),
	}

	return &out, nil
}
