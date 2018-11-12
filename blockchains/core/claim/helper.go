package claim

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/transfer"
)

func retrieveAllClaimsKeyname() string {
	return "claims"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Claim",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableClaim) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					str := fmt.Sprintf("the given storable Claim ID (%s) is invalid: %s", storable.ID, idErr.Error())
					return nil, errors.New(str)
				}

				transferID, transferIDErr := uuid.FromString(storable.TransferID)
				if transferIDErr != nil {
					str := fmt.Sprintf("the given storable Claim Transfer ID (%s) is invalid: %s", storable.DepositID, transferIDErr.Error())
					return nil, errors.New(str)
				}

				depositID, depositIDErr := uuid.FromString(storable.DepositID)
				if depositIDErr != nil {
					str := fmt.Sprintf("the given storable Claim Deposit ID (%s) is invalid: %s", storable.DepositID, depositIDErr.Error())
					return nil, errors.New(str)
				}

				// retrieve the transfer:
				transferMetaData := transfer.SDKFunc.CreateMetaData()
				transferIns, transferInsErr := rep.RetrieveByID(transferMetaData, &transferID)
				if transferInsErr != nil {
					return nil, transferInsErr
				}

				// retrieve the deposit:
				depositMetaData := deposit.SDKFunc.CreateMetaData()
				depositIns, depositInsErr := rep.RetrieveByID(depositMetaData, &depositID)
				if depositInsErr != nil {
					return nil, depositInsErr
				}

				if trans, ok := transferIns.(transfer.Transfer); ok {
					if dep, ok := depositIns.(deposit.Deposit); ok {
						out := createClaim(&id, trans, dep)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depositID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Transfer instance", transferID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableClaim); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableClaim)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableClaim),
	})
}
