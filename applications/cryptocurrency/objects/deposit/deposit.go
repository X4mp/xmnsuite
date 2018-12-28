package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type deposit struct {
	UUID *uuid.UUID      `json:"id"`
	Off  offer.Offer     `json:"offer"`
	Frm  address.Address `json:"from"`
	Am   int             `json:"amount"`
}

func createDeposit(id *uuid.UUID, off offer.Offer, frm address.Address, amount int) (Deposit, error) {
	out := deposit{
		UUID: id,
		Off:  off,
		Frm:  frm,
		Am:   amount,
	}

	return &out, nil
}

func createDepositFromNormalized(normalized *normalizedDeposit) (Deposit, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	offIns, offInsErr := offer.SDKFunc.CreateMetaData().Denormalize()(normalized.Offer)
	if offInsErr != nil {
		return nil, offInsErr
	}

	fromAddrIns, fromAddrInsErr := address.SDKFunc.CreateMetaData().Denormalize()(normalized.From)
	if fromAddrInsErr != nil {
		return nil, fromAddrInsErr
	}

	if off, ok := offIns.(offer.Offer); ok {
		if fromAddr, ok := fromAddrIns.(address.Address); ok {
			return createDeposit(&id, off, fromAddr, normalized.Amount)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", fromAddrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

func createDepositFromStorable(storable *storableDeposit, rep entity.Repository) (Deposit, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	offID, offIDErr := uuid.FromString(storable.OfferID)
	if offIDErr != nil {
		return nil, offIDErr
	}

	fromID, fromIDErr := uuid.FromString(storable.FromID)
	if fromIDErr != nil {
		return nil, fromIDErr
	}

	offIns, offInsErr := rep.RetrieveByID(offer.SDKFunc.CreateMetaData(), &offID)
	if offInsErr != nil {
		return nil, offInsErr
	}

	fromAddrIns, fromAddrInsErr := rep.RetrieveByID(address.SDKFunc.CreateMetaData(), &fromID)
	if fromAddrInsErr != nil {
		return nil, fromAddrInsErr
	}

	if off, ok := offIns.(offer.Offer); ok {
		if fromAddr, ok := fromAddrIns.(address.Address); ok {
			return createDeposit(&id, off, fromAddr, storable.Amount)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", fromAddrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *deposit) ID() *uuid.UUID {
	return obj.UUID
}

// Offer returns the offer
func (obj *deposit) Offer() offer.Offer {
	return obj.Off
}

// From returns the from address
func (obj *deposit) From() address.Address {
	return obj.Frm
}

// Amount returns the amount
func (obj *deposit) Amount() int {
	return obj.Am
}
