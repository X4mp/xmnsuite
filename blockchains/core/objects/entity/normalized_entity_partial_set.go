package entity

type normalizedPartialSet struct {
	Ins   []Normalized `json:"entities"`
	Idx   int          `json:"index"`
	TotAm int          `json:"total_amount"`
}

func createNormalizedPartialSet(ins PartialSet, metaData MetaData) (*normalizedPartialSet, error) {
	entities := ins.Instances()
	normalizedEntities := []Normalized{}
	for _, oneEntity := range entities {
		normalized, normalizedErr := metaData.Normalize()(oneEntity)
		if normalizedErr != nil {
			return nil, normalizedErr
		}

		normalizedEntities = append(normalizedEntities, normalized)
	}

	out := normalizedPartialSet{
		Ins:   normalizedEntities,
		Idx:   ins.Index(),
		TotAm: ins.TotalAmount(),
	}

	return &out, nil
}
