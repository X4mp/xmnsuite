package bank

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllBanksKeyname() string {
	return "banks"
}

func retrieveBankByCurrencyIDKeyname(currencyID *uuid.UUID) string {
	base := retrieveAllBanksKeyname()
	return fmt.Sprintf("%s:by_currency_id:%s", base, currencyID.String())
}

func retrieveBankByPledgeIDKeyname(pledge *uuid.UUID) string {
	base := retrieveAllBanksKeyname()
	return fmt.Sprintf("%s:by_pledge_id:%s", base, pledge.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Bank",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableBank); ok {
				return fromStorableToBank(storable, rep)
			}

			ptr := new(normalizedBank)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToBank(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if bnk, ok := ins.(Bank); ok {
				out, outErr := createNormalizedBank(bnk)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Bank instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedBank); ok {
				return fromNormalizedToBank(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Bank instance")
		},
		EmptyNormalized: new(normalizedBank),
		EmptyStorable:   new(storableBank),
	})
}
