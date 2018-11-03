package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

func retrieveAllWithdrawalsKeyname() string {
	return "withdrawals"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Withdrawal",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableWithdrawal) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				fromWalletID, fromWalletIDErr := uuid.FromString(storable.FromWalletID)
				if fromWalletIDErr != nil {
					return nil, fromWalletIDErr
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				walletIns, walletInsErr := rep.RetrieveByID(walletMetaData, &fromWalletID)
				if walletInsErr != nil {
					return nil, walletInsErr
				}

				if wal, ok := walletIns.(wallet.Wallet); ok {
					out := createWithdrawal(&id, wal, storable.Amount)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", fromWalletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableWithdrawal); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableWithdrawal)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableWithdrawal),
	})
}
