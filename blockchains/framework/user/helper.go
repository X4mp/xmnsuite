package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

func retrieveAllUserKeyname() string {
	return "users"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "User",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableUser) (entity.Entity, error) {
				// create the metadata:
				walletMetaData := wallet.SDKFunc.CreateMetaData()

				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
					PubKeyAsString: storable.PubKey,
				})

				walletID, walletIDErr := uuid.FromString(storable.WalletID)
				if walletIDErr != nil {
					return nil, walletIDErr
				}

				ins, insErr := rep.RetrieveByID(walletMetaData, &walletID)
				if insErr != nil {
					return nil, insErr
				}

				if wal, ok := ins.(wallet.Wallet); ok {
					out := createUser(&id, pubKey, storable.Shares, wal)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance and thererfore the given data cannot be transformed to a User instance", walletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableUser); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableUser)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableUser),
	})
}
