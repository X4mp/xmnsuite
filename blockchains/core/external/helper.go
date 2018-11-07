package external

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/link"
)

func retrieveAllExternalsKeyname() string {
	return "externals"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "EWallet",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableExternal) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				resourceID, resourceIDErr := uuid.FromString(storable.ResID)
				if resourceIDErr != nil {
					return nil, resourceIDErr
				}

				linkID, linkIDErr := uuid.FromString(storable.LinkID)
				if linkIDErr != nil {
					return nil, linkIDErr
				}

				// retrieve the link:
				linkMetaData := link.SDKFunc.CreateMetaData()
				linkIns, linkInsErr := rep.RetrieveByID(linkMetaData, &linkID)
				if linkInsErr != nil {
					return nil, linkInsErr
				}

				if lnk, ok := linkIns.(link.Link); ok {
					out := createExternal(&id, lnk, &resourceID)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Link instance", linkID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableExternal); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableExternal)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableExternal),
	})
}
