package developer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

func retrieveAllDevelopersKeyname() string {
	return "developers"
}

func retrieveDevelopersByUserIDKeyname(userID *uuid.UUID) string {
	return fmt.Sprintf("developers:by_user_id:%s", userID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Developer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableDeveloper) (entity.Entity, error) {
				// create the metadata:
				userMetaData := user.SDKFunc.CreateMetaData()

				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				userID, userIDErr := uuid.FromString(storable.UserID)
				if userIDErr != nil {
					return nil, userIDErr
				}

				ins, insErr := rep.RetrieveByID(userMetaData, &userID)
				if insErr != nil {
					return nil, insErr
				}

				if usr, ok := ins.(user.User); ok {
					out := createDeveloper(&id, usr, storable.Name, storable.Resume)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance and thererfore the given data cannot be transformed to a Developer instance", userID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableDeveloper); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedDeveloper)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createDeveloperFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if dev, ok := ins.(Developer); ok {
				return createNormalizedDeveloper(dev)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Developer instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedDeveloper); ok {
				return createDeveloperFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Developer instance")
		},
		EmptyStorable:   new(storableDeveloper),
		EmptyNormalized: new(normalizedDeveloper),
	})
}
