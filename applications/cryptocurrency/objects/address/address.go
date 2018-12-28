package address

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

type address struct {
	UUID *uuid.UUID    `json:"id"`
	Wal  wallet.Wallet `json:"wallet"`
	Addr string        `json:"address"`
}

func createAddress(id *uuid.UUID, wal wallet.Wallet, addr string) (Address, error) {
	out := address{
		UUID: id,
		Wal:  wal,
		Addr: addr,
	}

	return &out, nil
}

func createAddressFromNormalized(normalized *normalizedAddress) (Address, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	walIns, walInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Wallet)
	if walInsErr != nil {
		return nil, walInsErr
	}

	if wal, ok := walIns.(wallet.Wallet); ok {
		return createAddress(&id, wal, normalized.Address)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
}

func createAddressFromStorable(storable *storableAddress, rep entity.Repository) (Address, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	walletID, walletIDErr := uuid.FromString(storable.WalletID)
	if walletIDErr != nil {
		return nil, walletIDErr
	}

	walIns, walInsErr := rep.RetrieveByID(wallet.SDKFunc.CreateMetaData(), &walletID)
	if walInsErr != nil {
		return nil, walInsErr
	}

	if wal, ok := walIns.(wallet.Wallet); ok {
		return createAddress(&id, wal, storable.Address)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *address) ID() *uuid.UUID {
	return obj.UUID
}

// Wallet returns the wallet
func (obj *address) Wallet() wallet.Wallet {
	return obj.Wal
}

// Address returns the address
func (obj *address) Address() string {
	return obj.Addr
}
