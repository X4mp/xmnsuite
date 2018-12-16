package seed

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
)

func retrieveAllSeedsKeyname() string {
	return "seeds"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Seed",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableSeed); ok {

				lnkMetaData := link.SDKFunc.CreateMetaData()

				seedID, seedIDErr := uuid.FromString(storable.ID)
				if seedIDErr != nil {
					return nil, seedIDErr
				}

				lnkID, lnkIDErr := uuid.FromString(storable.LinkID)
				if lnkIDErr != nil {
					return nil, lnkIDErr
				}

				lnkIns, lnkInsErr := rep.RetrieveByID(lnkMetaData, &lnkID)
				if lnkInsErr != nil {
					return nil, lnkInsErr
				}

				if lnk, ok := lnkIns.(link.Link); ok {
					ip := net.ParseIP(storable.IP)
					out := createSeed(&seedID, lnk, ip, storable.Port)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Link instance", lnkIns.ID().String())
				return nil, errors.New(str)
			}

			ptr := new(normalizedSeed)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createSeedFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if seed, ok := ins.(Seed); ok {
				out, outErr := createNormalizedSeed(seed)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Seed instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedSeed); ok {
				return createSeedFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Seed instance")
		},
		EmptyNormalized: new(normalizedSeed),
		EmptyStorable:   new(storableSeed),
	})
}
