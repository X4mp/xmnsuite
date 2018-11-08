package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func keyname() string {
	return "genesis"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Genesis",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableGenesis) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				initialDepID, initialDepIDErr := uuid.FromString(storable.InitialDepositID)
				if initialDepIDErr != nil {
					return nil, initialDepIDErr
				}

				// retrieve the initial deposit:
				depositMet := deposit.SDKFunc.CreateMetaData()
				depIns, depInsErr := rep.RetrieveByID(depositMet, &initialDepID)
				if depInsErr != nil {
					return nil, depInsErr
				}

				if deposit, ok := depIns.(deposit.Deposit); ok {
					out := createGenesis(&id, storable.GzPricePerKb, storable.MxAmountOfValidators, deposit)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Genesis instance", initialDepID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableGenesis); ok {
				return fromStorableToEntity(storable)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedGenesis)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createGenesisFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Genesis instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if gen, ok := ins.(Genesis); ok {
				out, outErr := createNormalizedGenesis(gen)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedGenesis); ok {
				return createGenesisFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Genesis instance")
		},
		EmptyStorable: new(storableGenesis),
	})
}

func representation(depositRepresentation entity.Representation) entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if gen, ok := ins.(Genesis); ok {
				out := createStorableGenesis(gen)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			return []string{
				keyname(),
			}, nil
		},
		Sync: func(rep entity.Repository, service entity.Service, ins entity.Entity) error {
			if gen, ok := ins.(Genesis); ok {
				dep := gen.Deposit()
				metaData := depositRepresentation.MetaData()
				_, retDepErr := rep.RetrieveByID(metaData, dep.ID())
				if retDepErr != nil {
					saveErr := service.Save(dep, depositRepresentation)
					if saveErr != nil {
						return saveErr
					}
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
