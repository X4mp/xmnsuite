package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllDepositsKeyname() string {
	return "deposits"
}

func retrieveDepositsByBankIDKeyname(bankID *uuid.UUID) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_bank_id:%s", base, bankID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Deposit",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableDeposit); ok {
				return fromStorableToDeposit(storable, rep)
			}

			ptr := new(normalizedDeposit)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToDeposit(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if dep, ok := ins.(Deposit); ok {
				out, outErr := createNormalizedDeposit(dep)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedDeposit); ok {
				return fromNormalizedToDeposit(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Deposit instance")
		},
		EmptyNormalized: new(normalizedDeposit),
		EmptyStorable:   new(storableDeposit),
	})
}
