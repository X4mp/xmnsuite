package category

type normalizedCategory struct {
	ID          string     `json:"id"`
	Parent      Normalized `json:"parent"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

func createNormalizedCategory(ins Category) (*normalizedCategory, error) {
	var normalizedPar Normalized
	if ins.HasParent() {
		metaData := createMetaData()
		par, parErr := metaData.Normalize()(ins.Parent())
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
