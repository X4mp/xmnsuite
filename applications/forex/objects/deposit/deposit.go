package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type deposit struct {
	UUID *uuid.UUID `json:"id"`
	Am   int        `json:"amount"`
	Bnk  bank.Bank  `json:"bank"`
}

func createDeposit(id *uuid.UUID, amount int, bnk bank.Bank) (Deposit, error) {
	out := deposit{
		UUID: id,
		Am:   amount,
		Bnk:  bnk,
	}

	return &out, nil
}

func fromNormalizedToDeposit(normalized *normalizedDeposit) (Deposit, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	bnkIns, bnkInsErr := bank.SDKFunc.CreateMetaData().Denormalize()(normalized.Bank)
	if bnkInsErr != nil {
		return nil, bnkInsErr
	}

	if bnk, ok := bnkIns.(bank.Bank); ok {
		return createDeposit(&id, normalized.Amount, bnk)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Bank instance", bnkIns.ID().String())
	return nil, errors.New(str)
}

func fromStorableToDeposit(storable *storableDeposit, rep entity.Repository) (Deposit, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	bnkID, bnkIDErr := uuid.FromString(storable.BankID)
	if bnkIDErr != nil {
		return nil, bnkIDErr
	}

	bnkIns, bnkInsErr := rep.RetrieveByID(bank.SDKFunc.CreateMetaData(), &bnkID)
	if bnkInsErr != nil {
		return nil, bnkInsErr
	}

	if bnk, ok := bnkIns.(bank.Bank); ok {
		return createDeposit(&id, storable.Amount, bnk)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Bank instance", bnkIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *deposit) ID() *uuid.UUID {
	return obj.UUID
}

// Amount returns the amount
func (obj *deposit) Amount() int {
	return obj.Am
}

// Bank returns the bank
func (obj *deposit) Bank() bank.Bank {
	return obj.Bnk
}
