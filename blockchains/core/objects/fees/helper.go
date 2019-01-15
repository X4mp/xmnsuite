package fees

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllFeesKeyname() string {
	return "fees"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Fee",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableFee); ok {
				return createFeeFromStorable(storable, rep)
			}

			ptr := new(normalizedFee)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createFeeFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if req, ok := ins.(Fee); ok {
				return createNormalizedFee(req)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Fee instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedFee); ok {
				return createFeeFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Fee instance")
		},
		EmptyStorable:   new(storableFee),
		EmptyNormalized: new(normalizedFee),
	})
}
