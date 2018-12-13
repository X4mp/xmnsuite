package fiatchain

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllFiatChainsKeyname() string {
	return "fiatchains"
}

func retrieveFiatChainByDepositIDKeyname(depID *uuid.UUID) string {
	base := retrieveAllFiatChainsKeyname()
	return fmt.Sprintf("%s:by_deposit_id:%s", base, depID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "FiatChain",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableFiatChain); ok {
				return fromStorableToFiatChain(storable, rep)
			}

			ptr := new(normalizedFiatChain)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToFiatChain(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if fc, ok := ins.(FiatChain); ok {
				out, outErr := createNormalizedFiatChain(fc)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid FiatChain instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedFiatChain); ok {
				return fromNormalizedToFiatChain(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized FiatChain instance")
		},
		EmptyNormalized: new(normalizedFiatChain),
		EmptyStorable:   new(storableFiatChain),
	})
}
