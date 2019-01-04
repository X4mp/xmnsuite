package information

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

func keyname() string {
	return "information"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Information",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableInformation) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				return createInformation(&id, storable.ConcensusNeeded, storable.GzPricePerKb, storable.MxAmountOfValidators)
			}

			if storable, ok := data.(*storableInformation); ok {
				return fromStorableToEntity(storable)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedInformation)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createInformationFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Information instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if info, ok := ins.(Information); ok {
				out, outErr := createNormalizedInformation(info)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Information instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedInformation); ok {
				return createInformationFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Information instance")
		},
		EmptyStorable:   new(storableInformation),
		EmptyNormalized: new(normalizedInformation),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if info, ok := ins.(Information); ok {
				out := createStorableInformation(info)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Information instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			return []string{
				keyname(),
			}, nil
		},
		Sync: func(ds datastore.DataStore, ins entity.Entity) error {
			if info, ok := ins.(Information); ok {
				// crate metadata and representation:
				metaData := createMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				repository := createRepository(entityRepository, metaData)

				// retrieve the information:
				retInfo, retInfoErr := repository.Retrieve()
				if retInfoErr != nil {
					// the info does not exists, so return successfully:
					return nil
				}

				// make sure the infos have the same ids:
				if bytes.Compare(info.ID().Bytes(), retInfo.ID().Bytes()) != 0 {
					str := fmt.Sprintf("the given information (ID: %s) does not have the same ID as the currently saved info (ID: %s)", info.ID().String(), retInfo.ID().String())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Information instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
