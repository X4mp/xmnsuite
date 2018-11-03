package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/user"
)

func retrieveAllRequestsKeyname() string {
	return "requests"
}

func createMetaData(met entity.MetaData) entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Request",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableRequest) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				fromID, fromIDErr := uuid.FromString(storable.FromUserID)
				if fromIDErr != nil {
					return nil, fromIDErr
				}

				toEntity := met.ToEntity()
				newIns, newInsErr := toEntity(rep, storable.NewEntityJS)
				if newInsErr != nil {
					return nil, newInsErr
				}

				userMetaData := user.SDKFunc.CreateMetaData()
				from, fromErr := rep.RetrieveByID(userMetaData, &fromID)
				if fromErr != nil {
					return nil, fromErr
				}

				if fromUser, ok := from.(user.User); ok {
					out := createRequest(&id, fromUser, newIns)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid user instance", id.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableRequest); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableRequest)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableRequest),
	})
}
