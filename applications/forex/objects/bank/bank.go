package bank

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

type bank struct {
	UUID  *uuid.UUID        `json:"id"`
	Pldge pledge.Pledge     `json:"pledge"`
	Curr  currency.Currency `json:"currency"`
	Am    int               `json:"amount"`
	Pr    int               `json:"price"`
}

func createBank(id *uuid.UUID, pldge pledge.Pledge, curr currency.Currency, amount int, price int) (Bank, error) {
	out := bank{
		UUID:  id,
		Pldge: pldge,
		Curr:  curr,
		Am:    amount,
		Pr:    price,
	}

	return &out, nil
}

func fromNormalizedToBank(normalized *normalizedBank) (Bank, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	pldgeIns, pldgeInsErr := pledge.SDKFunc.CreateMetaData().Denormalize()(normalized.Pledge)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	currIns, currInsErr := currency.SDKFunc.CreateMetaData().Denormalize()(normalized.Currency)
	if currInsErr != nil {
		return nil, currInsErr
	}

	if pldge, ok := pldgeIns.(pledge.Pledge); ok {
		if curr, ok := currIns.(currency.Currency); ok {
			return createBank(&id, pldge, curr, normalized.Amount, normalized.Price)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Currency instance", currIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
	return nil, errors.New(str)
}

func fromStorableToBank(storable *storableBank, rep entity.Repository) (Bank, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	pldgeID, pldgeIDErr := uuid.FromString(storable.PledgeID)
	if pldgeIDErr != nil {
		return nil, pldgeIDErr
	}

	currID, currIDErr := uuid.FromString(storable.CurrencyID)
	if currIDErr != nil {
		return nil, currIDErr
	}

	pldgeIns, pldgeInsErr := rep.RetrieveByID(pledge.SDKFunc.CreateMetaData(), &pldgeID)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	currIns, currInsErr := rep.RetrieveByID(currency.SDKFunc.CreateMetaData(), &currID)
	if currInsErr != nil {
		return nil, currInsErr
	}

	if pldge, ok := pldgeIns.(pledge.Pledge); ok {
		if curr, ok := currIns.(currency.Currency); ok {
			return createBank(&id, pldge, curr, storable.Amount, storable.Price)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Currency instance", currIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *bank) ID() *uuid.UUID {
	return obj.UUID
}

// Pledge returns the pledge
func (obj *bank) Pledge() pledge.Pledge {
	return obj.Pldge
}

// Currency returns the currency
func (obj *bank) Currency() currency.Currency {
	return obj.Curr
}

// Amount returns the amount
func (obj *bank) Amount() int {
	return obj.Am
}

// Price returns the price
func (obj *bank) Price() int {
	return obj.Pr
}
