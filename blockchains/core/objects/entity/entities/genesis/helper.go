package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/datastore"
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

				userID, userIDErr := uuid.FromString(storable.UserID)
				if userIDErr != nil {
					return nil, userIDErr
				}

				infID, infIDErr := uuid.FromString(storable.InfoID)
				if infIDErr != nil {
					return nil, infIDErr
				}

				// retrieve the initial deposit:
				depositMet := deposit.SDKFunc.CreateMetaData()
				depIns, depInsErr := rep.RetrieveByID(depositMet, &initialDepID)
				if depInsErr != nil {
					return nil, depInsErr
				}

				// retrieve the user:
				usrMet := user.SDKFunc.CreateMetaData()
				usrIns, usrInsErr := rep.RetrieveByID(usrMet, &userID)
				if usrInsErr != nil {
					return nil, usrInsErr
				}

				// retrieve the information:
				infMet := information.SDKFunc.CreateMetaData()
				infIns, infInsErr := rep.RetrieveByID(infMet, &infID)
				if infInsErr != nil {
					return nil, infInsErr
				}

				if deposit, ok := depIns.(deposit.Deposit); ok {
					if usr, ok := usrIns.(user.User); ok {
						if inf, ok := infIns.(information.Information); ok {
							return createGenesis(&id, inf, deposit, usr)
						}

						str := fmt.Sprintf("the entity (ID: %s) is not a valid Information instance", infID.String())
						return nil, errors.New(str)
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", userID.String())
					return nil, errors.New(str)
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
		EmptyStorable:   new(storableGenesis),
		EmptyNormalized: new(normalizedGenesis),
	})
}

func representation() entity.Representation {
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
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if gen, ok := ins.(Genesis); ok {
				// the user must have enough shares in order to fill the concensus, on genesis:
				genUser := gen.User()
				genTo := gen.Deposit().To()
				if genUser.Shares() < genTo.ConcensusNeeded() {
					str := fmt.Sprintf("the genesis user (ID: %s) does not have enough shares (%d) in order to convert the concensus of the genesis wallet (ID: %s ConcensusNeeded: %d)", genUser.ID().String(), genUser.Shares(), genTo.ID(), genTo.ConcensusNeeded())
					return errors.New(str)
				}

				// crate metadata and representation:
				metaData := createMetaData()

				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)
				genRepository := createRepository(repository, metaData)

				// if the genesis exists, return error:
				_, retGenErr := genRepository.Retrieve()
				if retGenErr == nil {
					str := fmt.Sprintf("there is already a Genesis instance")
					return errors.New(str)
				}

				// deposit:
				dep := gen.Deposit()
				depRepresentation := deposit.SDKFunc.CreateRepresentation()
				_, retDepErr := repository.RetrieveByID(depRepresentation.MetaData(), dep.ID())
				if retDepErr == nil {
					str := fmt.Sprintf("the Genesis instance contains a Deposit instance (ID: %s) that is already saved", dep.ID().String())
					return errors.New(str)
				}

				depSaveErr := service.Save(dep, depRepresentation)
				if depSaveErr != nil {
					return depSaveErr
				}

				// user:
				usr := gen.User()
				usrRepresentation := user.SDKFunc.CreateRepresentation()
				_, retUsrErr := repository.RetrieveByID(usrRepresentation.MetaData(), usr.ID())
				if retUsrErr == nil {
					str := fmt.Sprintf("the Genesis instance contains a User instance (ID: %s) that is already saved", usr.ID().String())
					return errors.New(str)
				}

				usrSaveErr := service.Save(usr, usrRepresentation)
				if usrSaveErr != nil {
					return usrSaveErr
				}

				// information:
				inf := gen.Info()
				infoRepresentation := information.SDKFunc.CreateRepresentation()
				_, retInfErr := repository.RetrieveByID(infoRepresentation.MetaData(), inf.ID())
				if retInfErr == nil {
					str := fmt.Sprintf("the Genesis instance contains an Information instance (ID: %s) that is already saved", inf.ID().String())
					return errors.New(str)
				}

				infSaveErr := service.Save(inf, infoRepresentation)
				if infSaveErr != nil {
					return infSaveErr
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
