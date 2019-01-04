package wallet

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
)

func retrieveAllWalletKeyname() string {
	return "wallets"
}

func retrieveByPublicKeyWalletKeyname(pubKey crypto.PublicKey) string {
	base := retrieveAllWalletKeyname()
	return fmt.Sprintf("%s:by_pubkey:%s", base, pubKey.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Wallet",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableWallet); ok {
				return createWalletFromStorable(storable)
			}

			ptr := new(storableWallet)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createWalletFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if wallet, ok := ins.(Wallet); ok {
				out := createStoredWallet(wallet)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if storable, ok := ins.(*storableWallet); ok {
				return createWalletFromStorable(storable)
			}

			return nil, errors.New("the given instance is not a valid normalized Wallet instance")
		},
		EmptyStorable:   new(storableWallet),
		EmptyNormalized: new(storableWallet),
	})
}
