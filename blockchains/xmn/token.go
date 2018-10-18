package xmn

import (
	"errors"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type token struct {
	Sym  string `json:"symbol"`
	Nme  string `json:"name"`
	Desc string `json:"description"`
}

func createToken(symbol string, name string, desc string) Token {
	out := token{
		Sym:  symbol,
		Nme:  name,
		Desc: desc,
	}

	return &out
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

type tokenService struct {
	keyname string
	store   datastore.DataStore
}

func createTokenService(store datastore.DataStore) TokenService {
	out := tokenService{
		keyname: "token",
		store:   store,
	}

	return &out
}

// Retrieve retrieves the token
func (app *tokenService) Retrieve() (Token, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: app.keyname,
		Obj: new(token),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		return nil, errors.New("there was an error while retrieving the Token instance")
	}

	// cast the instance:
	if tok, ok := obj.Obj.(Token); ok {
		return tok, nil
	}

	return nil, errors.New("the retrieved data cannot be casted to a Token instance")
}

// Save saves the token
func (app *tokenService) Save(tok Token) error {
	// make sure the instance does not exists already:
	_, retErr := app.Retrieve()
	if retErr == nil {
		return errors.New("the Token instance already exists")
	}

	// save the object:
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: app.keyname,
		Obj: tok,
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the Token instance")
	}

	return nil
}
