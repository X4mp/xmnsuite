package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func keyname() string {
	return "genesis"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Genesis",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableGenesis) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				initialDepID, initialDepIDErr := uuid.FromString(storable.InitialDepositID)
				if initialDepIDErr != nil {
					return nil, initialDepIDErr
				}

				// retrieve the initial deposit:
				depositMet := deposit.SDKFunc.CreateMetaData()
				depIns, depInsErr := rep.RetrieveByID(depositMet, &initialDepID)
				if depInsErr != nil {
					return nil, depInsErr
				}

				if deposit, ok := depIns.(deposit.Deposit); ok {
					out := createGenesis(&id, storable.GzPricePerKb, storable.MxAmountOfValidators, deposit)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid InitialDeposit instance", initialDepID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableGenesis); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableGenesis)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableGenesis),
	})
}
