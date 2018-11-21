package deposit

import (
	"errors"
	"fmt"
	"math"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
)

type deposit struct {
	UUID     *uuid.UUID    `json:"id"`
	ToWallet wallet.Wallet `json:"to"`
	Tok      token.Token   `json:"token"`
	Am       int           `json:"amount"`
}

func createDeposit(id *uuid.UUID, toWallet wallet.Wallet, tok token.Token, amount int) (Deposit, error) {

	if amount <= 0 {
		return nil, errors.New("the amount (%d) must be bigger than 0")
	}

	if amount > math.MaxInt64-1 {
		str := fmt.Sprintf("the amount (%d) cannot be bigger than %d", amount, math.MaxInt64-1)
		return nil, errors.New(str)
	}

	out := deposit{
		UUID:     id,
		ToWallet: toWallet,
		Tok:      tok,
		Am:       amount,
	}

	return &out, nil
}

func createDepositFromNormalized(ins *normalizedDeposit) (Deposit, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		str := fmt.Sprintf("the given storable Deposit ID (%s) is invalid: %s", ins.ID, idErr.Error())
		return nil, errors.New(str)
	}

	toWalletIns, toWalletInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(ins.To)
	if toWalletInsErr != nil {
		return nil, toWalletInsErr
	}

	tokenIns, tokenInsErr := token.SDKFunc.CreateMetaData().Denormalize()(ins.Token)
	if tokenInsErr != nil {
		return nil, tokenInsErr
	}

	if toWallet, ok := toWalletIns.(wallet.Wallet); ok {
		if tok, ok := tokenIns.(token.Token); ok {
			return createDeposit(&id, toWallet, tok, ins.Amount)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not  avalid Token instance", tokenIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not  avalid Wallet instance", toWalletIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (app *deposit) ID() *uuid.UUID {
	return app.UUID
}

// To returns the to user
func (app *deposit) To() wallet.Wallet {
	return app.ToWallet
}

// Token returns the token
func (app *deposit) Token() token.Token {
	return app.Tok
}

// Amount returns the amount
func (app *deposit) Amount() int {
	return app.Am
}
