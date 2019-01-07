package keyname

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
)

func retrieveAllKeynamesKeyname() string {
	return "keynames"
}

func retrieveKeynameByNameKeyname(name string) string {
	base := retrieveAllKeynamesKeyname()
	return fmt.Sprintf("%s:by_name:%s", base, name)
}

func retrieveKeynameByGroupKeyname(grp group.Group) string {
	base := retrieveAllKeynamesKeyname()
	return fmt.Sprintf("%s:by_group_id:%s", base, grp.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Keyname",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableKeyname); ok {
				return createKeynameFromStorable(rep, storable)
			}

			ptr := new(normalizedKeyname)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createKeynameFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if kname, ok := ins.(Keyname); ok {
				return createNormalizedKeyname(kname)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Keyname instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedKeyname); ok {
				return createKeynameFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Keyname instance")
		},
		EmptyStorable:   new(storableKeyname),
		EmptyNormalized: new(normalizedKeyname),
	})
}
