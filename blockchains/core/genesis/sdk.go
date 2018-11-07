package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Genesis represents the genesis instance
type Genesis interface {
	ID() *uuid.UUID
	GazPricePerKb() int
	MaxAmountOfValidators() int
	Deposit() deposit.Deposit
}

// Service represents the Genesis service
type Service interface {
	Save(ins Genesis) error
}

// Repository represents the Genesis repository
type Repository interface {
	Retrieve() (Genesis, error)
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	DepositRepresentation entity.Representation
}

// SDKFunc represents the Genesis SDK func
var SDKFunc = struct {
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
	CreateRepository: func(params CreateRepositoryParams) Repository {
		met := createMetaData()
		out := createRepository(params.EntityRepository, met)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
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
				saveIfNotExists := func(representation entity.Representation, ins entity.Entity) error {
					metaData := representation.MetaData()
					_, retDepErr := rep.RetrieveByID(metaData, ins.ID())
					if retDepErr != nil {
						saveErr := service.Save(ins, representation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				if gen, ok := ins.(Genesis); ok {
					saveIfNotExists(params.DepositRepresentation, gen.Deposit())
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
