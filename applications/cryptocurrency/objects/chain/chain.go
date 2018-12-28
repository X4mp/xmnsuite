package chain

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type chain struct {
	UUID  *uuid.UUID  `json:"id"`
	Off   offer.Offer `json:"offer"`
	TotAm int         `json:"total_amount"`
}

func createChain(id *uuid.UUID, off offer.Offer, totalAmount int) (Chain, error) {
	out := chain{
		UUID:  id,
		Off:   off,
		TotAm: totalAmount,
	}

	return &out, nil
}

func createChainFromNormalized(normalized *normalizedChain) (Chain, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	offIns, offInsErr := offer.SDKFunc.CreateMetaData().Denormalize()(normalized.Offer)
	if offInsErr != nil {
		return nil, offInsErr
	}

	if off, ok := offIns.(offer.Offer); ok {
		return createChain(&id, off, normalized.TotalAmount)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

func createChainFromStorable(storable *storableChain, rep entity.Repository) (Chain, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	offerID, offerIDErr := uuid.FromString(storable.OfferID)
	if offerIDErr != nil {
		return nil, offerIDErr
	}

	offIns, offInsErr := rep.RetrieveByID(offer.SDKFunc.CreateMetaData(), &offerID)
	if offInsErr != nil {
		return nil, offInsErr
	}

	if off, ok := offIns.(offer.Offer); ok {
		return createChain(&id, off, storable.TotalAmount)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *chain) ID() *uuid.UUID {
	return obj.UUID
}

// Offer returns the offer
func (obj *chain) Offer() offer.Offer {
	return obj.Off
}

// TotalAmount returns the total amount
func (obj *chain) TotalAmount() int {
	return obj.TotAm
}
