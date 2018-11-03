package genesis

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
)

/*
 * Genesis
 */

type genesis struct {
	UUID                 *uuid.UUID      `json:"id"`
	GzPricePerKb         int             `json:"gaz_price_per_kb"`
	MxAmountOfValidators int             `json:"max_amount_of_validators"`
	Dep                  deposit.Deposit `json:"deposit"`
	Tok                  token.Token     `json:"token"`
}

func createGenesis(id *uuid.UUID, gazPricePerKb int, maxAmountOfValidators int, dep deposit.Deposit, tok token.Token) Genesis {
	out := genesis{
		UUID:                 id,
		GzPricePerKb:         gazPricePerKb,
		MxAmountOfValidators: maxAmountOfValidators,
		Dep:                  dep,
		Tok:                  tok,
	}

	return &out
}

// ID returns the ID
func (app *genesis) ID() *uuid.UUID {
	return app.UUID
}

// GazPricePerKb returns the gazPricePerKb
func (app *genesis) GazPricePerKb() int {
	return app.GzPricePerKb
}

// MaxAmountOfValidators returns the maxAmountOfValidators
func (app *genesis) MaxAmountOfValidators() int {
	return app.MxAmountOfValidators
}

// Deposit returns the initial deposit
func (app *genesis) Deposit() deposit.Deposit {
	return app.Dep
}

// Token returns the token
func (app *genesis) Token() token.Token {
	return app.Tok
}
