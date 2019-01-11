package category

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllCategoriesKeyname() string {
	return "categories"
}

func retrieveCategoryByParentCategoryKeyname(parent Category) string {
	base := retrieveAllCategoriesKeyname()
	return fmt.Sprintf("%s:by_parent_category_id:%s", base, parent.ID().String())
}

func retrieveCategoryWithoutParentKeyname() string {
	base := retrieveAllCategoriesKeyname()
	return fmt.Sprintf("%s:without_parent", base)
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Category",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableCategory); ok {
				return createCategoryFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedCategory)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createCategoryFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Category instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if cat, ok := ins.(Category); ok {
				out, outErr := createNormalizedCategory(cat)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedCategory); ok {
				return createCategoryFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Category instance")
		},
		EmptyStorable:   new(storableCategory),
		EmptyNormalized: new(normalizedCategory),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if cat, ok := ins.(Category); ok {
				out := createStorableCategory(cat)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if cat, ok := ins.(Category); ok {
				keynames := []string{
					retrieveAllCategoriesKeyname(),
				}

				if cat.HasParent() {
					keynames = append(keynames, retrieveCategoryByParentCategoryKeyname(cat.Parent()))
				}

				if !cat.HasParent() {
					keynames = append(keynames, retrieveCategoryWithoutParentKeyname())
				}

				return keynames, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if cat, ok := ins.(Category); ok {
				// crate metadata and representation:
				rep := representation()
				metaData := rep.MetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				entityService := entity.SDKFunc.CreateService(ds)

				// if the parent category does not exists, save it:
				if cat.HasParent() {
					_, parCatErr := entityRepository.RetrieveByID(metaData, cat.Parent().ID())
					if parCatErr == nil {
						return nil
					}

					// save the parent category:
					saveParentErr := entityService.Save(cat.Parent(), rep)
					if saveParentErr != nil {
						return saveParentErr
					}
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
