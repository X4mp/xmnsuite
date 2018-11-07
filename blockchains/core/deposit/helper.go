package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

func retrieveAllDepositsKeyname() string {
	return "deposits"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Deposit",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableDeposit) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				toWalletID, toWalletIDErr := uuid.FromString(storable.ToWalletID)
				if toWalletIDErr != nil {
					return nil, toWalletIDErr
				}

				tokenID, tokenIDErr := uuid.FromString(storable.TokenID)
				if tokenIDErr != nil {
					return nil, tokenIDErr
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				walletIns, walletInsErr := rep.RetrieveByID(walletMetaData, &toWalletID)
				if walletInsErr != nil {
					return nil, walletInsErr
				}

				// retrieve the token:
				tokenMetaData := token.SDKFunc.CreateMetaData()
				tokenIns, tokenInsErr := rep.RetrieveByID(tokenMetaData, &tokenID)
				if tokenInsErr != nil {
					return nil, tokenInsErr
				}

				if wal, ok := walletIns.(wallet.Wallet); ok {
					if tok, ok := tokenIns.(token.Token); ok {
						out := createDeposit(&id, wal, tok, storable.Amount)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Token instance", tokenID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", toWalletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableDeposit); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableDeposit)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableDeposit),
	})
}
