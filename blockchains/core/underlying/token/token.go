package token

import (
	uuid "github.com/satori/go.uuid"
)

type token struct {
	UUID *uuid.UUID `json:"id"`
	Sym  string     `json:"symbol"`
	Nme  string     `json:"name"`
	Desc string     `json:"description"`
}

func createToken(id *uuid.UUID, symbol string, name string, desc string) Token {
	out := token{
		UUID: id,
		Sym:  symbol,
		Nme:  name,
		Desc: desc,
	}

	return &out
}

func createTokenFromStorable(ins *storableToken) (Token, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	out := createToken(&id, ins.Symbol, ins.Name, ins.Description)
	return out, nil
}

// ID returns the ID
func (app *token) ID() *uuid.UUID {
	return app.UUID
}

// Symbol returns the symbol
func (app *token) Symbol() string {
	return app.Sym
}

// Name returns the name
func (app *token) Name() string {
	return app.Nme
}

// Description returns the description
func (app *token) Description() string {
	return app.Desc
}
