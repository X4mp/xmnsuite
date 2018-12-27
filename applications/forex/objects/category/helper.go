package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllCategoriesKeyname() string {
	return "categories"
}

func retrieveCategoriesByParentCategoryIDKeyname(parentID *uuid.UUID) string {
	base := retrieveAllCategoriesKeyname()
	return fmt.Sprintf("%s:by_parent_id:%s", base, parentID.String())
}

func retrieveCcategoriesWithoutParentKeyname() string {
	base := retrieveAllCategoriesKeyname()
	return fmt.Sprintf("%s:without_parent", base)
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

func toData(cat Category) *Data {
	var parent *Data
	if cat.HasParent() {
		parent = toData(cat.Parent())
	}

	out := Data{
		ID:          cat.ID().String(),
		Parent:      parent,
		Name:        cat.Name(),
		Description: cat.Description(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	categories := []*Data{}
	for _, oneIns := range ins {
		if cat, ok := oneIns.(Category); ok {
			categories = append(categories, toData(cat))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Category instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Categories:  categories,
	}

	return &out, nil
}
