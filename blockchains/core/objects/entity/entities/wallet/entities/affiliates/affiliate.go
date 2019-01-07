package affiliates

import (
	"errors"
	"fmt"
	"net/url"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type affiliate struct {
	UUID *uuid.UUID    `json:"id"`
	Own  wallet.Wallet `json:"owner"`
	UR   *url.URL      `json:"url"`
}

func createAffiliate(id *uuid.UUID, owner wallet.Wallet, url *url.URL) (Affiliate, error) {
	out := affiliate{
		UUID: id,
		Own:  owner,
		UR:   url,
	}

	return &out, nil
}

func createAffiliateFromNormalized(normalized *normalizedAffiliate) (Affiliate, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	ownIns, ownInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Owner)
	if ownInsErr != nil {
		return nil, ownInsErr
	}

	ur, urErr := url.Parse(normalized.URL)
	if urErr != nil {
		return nil, urErr
	}

	if own, ok := ownIns.(wallet.Wallet); ok {
		return createAffiliate(&id, own, ur)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", ownIns.ID().String())
	return nil, errors.New(str)
}

func createAffiliateFromStorable(storable *storableAffiliate, rep entity.Repository) (Affiliate, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	walletID, walletIDErr := uuid.FromString(storable.OwnerID)
	if walletIDErr != nil {
		return nil, walletIDErr
	}

	ownIns, ownInsErr := rep.RetrieveByID(wallet.SDKFunc.CreateMetaData(), &walletID)
	if ownInsErr != nil {
		return nil, ownInsErr
	}

	ur, urErr := url.Parse(storable.URL)
	if urErr != nil {
		return nil, urErr
	}

	if own, ok := ownIns.(wallet.Wallet); ok {
		return createAffiliate(&id, own, ur)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", ownIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *affiliate) ID() *uuid.UUID {
	return obj.UUID
}

// Owner returns the owner
func (obj *affiliate) Owner() wallet.Wallet {
	return obj.Own
}

// URL returns the url
func (obj *affiliate) URL() *url.URL {
	return obj.UR
}
