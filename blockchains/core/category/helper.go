package category

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllCategoriesKeyname() string {
	return "categories"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Category",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableCategory); ok {
				return createCategoryFromStorable(storable)
			}

			ptr := new(storableCategory)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createCategoryFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if cat, ok := ins.(Category); ok {
				return createStorableCategory(cat), nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*storableCategory); ok {
				return createCategoryFromStorable(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Category instance")
		},
		EmptyStorable: new(storableCategory),
	})
}
