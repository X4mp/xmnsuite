package currency

type storableCurrency struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createStorableCurrency(ins Currency) *storableCurrency {
	out := storableCurrency{
		ID:          ins.ID().String(),
		CategoryID:  ins.Category().ID().String(),
		Symbol:      ins.Symbol(),
		Name:        ins.Name(),
		Description: ins.Description(),
	}

	return &out
}
