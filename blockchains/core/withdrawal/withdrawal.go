package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

type withdrawal struct {
	UUID       *uuid.UUID    `json:"id"`
	FromWallet wallet.Wallet `json:"from"`
	Tok        token.Token   `json:"token"`
	Am         int           `json:"amount"`
}

func createWithdrawal(id *uuid.UUID, fromWallet wallet.Wallet, tok token.Token, amount int) Withdrawal {
	out := withdrawal{
		UUID:       id,
		FromWallet: fromWallet,
		Tok:        tok,
		Am:         amount,
	}

	return &out
}

func createWithdrawalFromNormalized(ins *normalizedWithdrawal) (Withdrawal, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		str := fmt.Sprintf("the given storable Withdrawal ID (%s) is invalid: %s", ins.ID, idErr.Error())
		return nil, errors.New(str)
	}

	fromIns, fromInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(ins.From)
	if fromInsErr != nil {
		return nil, fromInsErr
	}

	tokIns, tokInsErr := token.SDKFunc.CreateMetaData().Denormalize()(ins.Token)
	if tokInsErr != nil {
		return nil, tokInsErr
	}

	if from, ok := fromIns.(wallet.Wallet); ok {
		if tok, ok := tokIns.(token.Token); ok {
			out := createWithdrawal(&id, from, tok, ins.Amount)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Token instance", tokIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", fromIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (app *withdrawal) ID() *uuid.UUID {
	return app.UUID
}

// From returns the from wallet
func (app *withdrawal) From() wallet.Wallet {
	return app.FromWallet
}

// Token returns the token
func (app *withdrawal) Token() token.Token {
	return app.Tok
}

// Amount returns the amount
func (app *withdrawal) Amount() int {
	return app.Am
}
