package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

func retrieveAllPledgesKeyname() string {
	return "pledges"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Pledge",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storablePledge) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				withdrawalID, withdrawalIDErr := uuid.FromString(storable.FromWithdrawalID)
				if withdrawalIDErr != nil {
					return nil, withdrawalIDErr
				}

				walletID, walletIDErr := uuid.FromString(storable.ToWalletID)
				if walletIDErr != nil {
					return nil, walletIDErr
				}

				// retrieve the withdrawal:
				withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
				fromIns, fromInsErr := rep.RetrieveByID(withdrawalMetaData, &withdrawalID)
				if fromInsErr != nil {
					return nil, fromInsErr
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				toIns, toInsErr := rep.RetrieveByID(walletMetaData, &walletID)
				if toInsErr != nil {
					return nil, toInsErr
				}

				if from, ok := fromIns.(withdrawal.Withdrawal); ok {
					if to, ok := toIns.(wallet.Wallet); ok {
						out := createPledge(&id, from, to, storable.Amount)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walletID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withdrawalID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storablePledge); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storablePledge)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storablePledge),
	})
}
