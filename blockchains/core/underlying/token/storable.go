package token

type storableToken struct {
	ID          string `json:"id"`
	Symbol      string `json:"string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createStorableToken(tok Token) *storableToken {
	out := storableToken{
		ID:          tok.ID().String(),
		Symbol:      tok.Symbol(),
		Name:        tok.Name(),
		Description: tok.Description(),
	}

	return &out
}
