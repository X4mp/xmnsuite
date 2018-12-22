package group

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllGroupsKeyname() string {
	return "groups"
}

func retrieveGroupByNameKeyname(name string) string {
	base := retrieveAllGroupsKeyname()
	return fmt.Sprintf("%s:by_name:%s", base, name)
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Group",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableGroup); ok {
				return createGroupFromStorable(storable)
			}

			ptr := new(storableGroup)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createGroupFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if grp, ok := ins.(Group); ok {
				out := createStorableGroup(grp)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Group instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*storableGroup); ok {
				return createGroupFromStorable(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Group instance")
		},
		EmptyStorable:   new(storableGroup),
		EmptyNormalized: new(storableGroup),
	})
}
