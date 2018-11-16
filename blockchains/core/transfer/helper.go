package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
)

func retrieveAllTransfersKeyname() string {
	return "transfers"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Transfer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableTransfer) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				withdrawalID, withdrawalIDErr := uuid.FromString(storable.WithdrawalID)
				if withdrawalIDErr != nil {
					return nil, withdrawalIDErr
				}

				depositID, depositIDErr := uuid.FromString(storable.DepositID)
				if depositIDErr != nil {
					return nil, depositIDErr
				}

				// retrieve the withdrawal:
				withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
				withdrawalIns, withdrawalInsErr := rep.RetrieveByID(withdrawalMetaData, &withdrawalID)
				if withdrawalInsErr != nil {
					return nil, withdrawalInsErr
				}

				// retrieve the deposit:
				depositMetaData := deposit.SDKFunc.CreateMetaData()
				depositIns, depositInsErr := rep.RetrieveByID(depositMetaData, &depositID)
				if depositInsErr != nil {
					return nil, depositInsErr
				}

				if withdrawl, ok := withdrawalIns.(withdrawal.Withdrawal); ok {
					if dep, ok := depositIns.(deposit.Deposit); ok {
						out := createTransfer(&id, withdrawl, dep)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depositID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withdrawalID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableTransfer); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableTransfer)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableTransfer),
	})
}
