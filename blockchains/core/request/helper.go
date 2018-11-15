package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

func retrieveAllRequestsKeyname() string {
	return "requests"
}

func createMetaData(reg Registry) entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Request",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableRequest) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					str := fmt.Sprintf("the storable ID (%s) is invalid: %s", storable.ID, idErr.Error())
					return nil, errors.New(str)
				}

				fromID, fromIDErr := uuid.FromString(storable.FromUserID)
				if fromIDErr != nil {
					str := fmt.Sprintf("the storable FromUserID (%s) is invalid: %s", storable.FromUserID, fromIDErr.Error())
					return nil, errors.New(str)
				}

				newIns, newInsErr := reg.FromJSONToEntity(storable.NewEntityJS)
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

				str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", id.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableRequest); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedRequest)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createRequestFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if req, ok := ins.(Request); ok {
				return createNormalizedRequest(req)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedRequest); ok {
				return createRequestFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Request instance")
		},
		EmptyStorable:   new(storableRequest),
		EmptyNormalized: new(normalizedRequest),
	})
}
