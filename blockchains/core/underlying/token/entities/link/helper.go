package link

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllLinksKeyname() string {
	return "links"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Link",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableLink) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				return createLink(&id, storable.Title, storable.Description)
			}

			if storable, ok := data.(*storableLink); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedLink)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createLinkFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if lnk, ok := ins.(Link); ok {
				return createNormalizedLink(lnk)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Link instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedLink); ok {
				return createLinkFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Link instance")
		},
		EmptyNormalized: new(normalizedLink),
		EmptyStorable:   new(storableLink),
	})
}
