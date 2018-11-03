package genesis

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

// Genesis represents the genesis instance
type Genesis interface {
	ID() *uuid.UUID
	GazPricePerKb() int
	MaxAmountOfValidators() int
	Deposit() deposit.Deposit
	Token() token.Token
}

// Service represents the Genesis service
type Service interface {
	Save(ins Genesis) error
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	InitialDepositMetaData       entity.MetaData
	InitialDepositRepresentation entity.Representation
	TokenMetaData                entity.MetaData
	TokenRepresentation          entity.Representation
}

// SDKFunc represents the Genesis SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
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
				saveIfNotExists := func(metaData entity.MetaData, representation entity.Representation, ins entity.Entity) error {
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
					saveIfNotExists(params.InitialDepositMetaData, params.InitialDepositRepresentation, gen.Deposit())
					saveIfNotExists(params.TokenMetaData, params.TokenRepresentation, gen.Token())
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Genesis instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
