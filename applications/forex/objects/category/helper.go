package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllCurrenciesKeyname() string {
	return "currencies"
}

func retrieveCurrenciesByParentCategoryIDKeyname(parentID *uuid.UUID) string {
	base := retrieveAllCurrenciesKeyname()
	return fmt.Sprintf("%s:by_parent_id:%s", base, parentID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Category",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableCategory); ok {
				return fromStorableToCategory(storable, rep)
			}

			ptr := new(normalizedCategory)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToCategory(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if curr, ok := ins.(Category); ok {
				return createNormalizedCategory(curr)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedCategory); ok {
				return fromNormalizedToCategory(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Category instance")
		},
		EmptyNormalized: new(normalizedCategory),
		EmptyStorable:   new(storableCategory),
	})
}
