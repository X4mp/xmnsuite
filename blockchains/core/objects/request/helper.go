package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

func retrieveAllRequestsKeyname() string {
	return "requests"
}

func createMetaData(reg *registry) entity.MetaData {
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

				knameID, knameIDErr := uuid.FromString(storable.KeynameID)
				if knameIDErr != nil {
					str := fmt.Sprintf("the storable KeynameID (%s) is invalid: %s", storable.KeynameID, knameIDErr.Error())
					return nil, errors.New(str)
				}

				fromIns, fromInsErr := rep.RetrieveByID(user.SDKFunc.CreateMetaData(), &fromID)
				if fromInsErr != nil {
					return nil, fromInsErr
				}

				knameIns, knameInsErr := rep.RetrieveByID(keyname.SDKFunc.CreateMetaData(), &knameID)
				if knameInsErr != nil {
					return nil, knameInsErr
				}

				if fromUser, ok := fromIns.(user.User); ok {
					if kname, ok := knameIns.(keyname.Keyname); ok {

						if storable.SaveEntityJSON != nil {
							saveIns, saveInsErr := reg.fromJSONToEntity(storable.SaveEntityJSON, kname.Name())
							if saveInsErr != nil {
								return nil, saveInsErr
							}

							out := createRequestWithSaveEntity(&id, fromUser, saveIns, storable.Reason, kname)
							return out, nil
						}

						delIns, delInsErr := reg.fromJSONToEntity(storable.DeleteEntityJSON, kname.Name())
						if delInsErr != nil {
							return nil, delInsErr
						}

						out := createRequestWithDeleteEntity(&id, fromUser, delIns, storable.Reason, kname)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Keyname instance", knameIns.ID().String())
					return nil, errors.New(str)
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
