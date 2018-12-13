package currency

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllCurrenciesKeyname() string {
	return "currencies"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Currency",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableCurrency); ok {
				return fromStorableToCurrency(storable, rep)
			}

			ptr := new(normalizedCurrency)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToCurrency(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if curr, ok := ins.(Currency); ok {
				return createNormalizedCurrency(curr)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Currency instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedCurrency); ok {
				return fromNormalizedToCurrency(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Currency instance")
		},
		EmptyNormalized: new(normalizedCurrency),
		EmptyStorable:   new(storableCurrency),
	})
}
