package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet/request/entities/pledge"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Validator represents a validator
type Validator interface {
	ID() *uuid.UUID
	PubKey() tcrypto.PubKey
	Pledge() pledge.Pledge
}

// Normalized represents the normalized validator
type Normalized interface {
}

// Repository represents the validator repository
type Repository interface {
	RetrieveSet(amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	PubKey tcrypto.PubKey
	Pledge pledge.Pledge
}

// SDKFunc represents the Validator SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Validator
	CreateRepository     func(ds datastore.DataStore) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Validator {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createValidator(params.ID, params.PubKey, params.Pledge)
		return out
	},
	CreateRepository: func(ds datastore.DataStore) Repository {
		met := createMetaData()
		entityRepository := entity.SDKFunc.CreateRepository(ds)
		out := createRepository(entityRepository, met)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if validator, ok := ins.(Validator); ok {
					out := createStorableValidator(validator)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Validator instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllValidatorsKeyname(),
				}, nil
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				// create the representations:
				pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()

				if val, ok := ins.(Validator); ok {
					// try to retrieve the pledge, save it if it doesnt exists:
					pldge := val.Pledge()
					_, retPledgeErr := repository.RetrieveByID(pledgeRepresentation.MetaData(), pldge.ID())
					if retPledgeErr != nil {
						// save the pledge:
						saveErr := service.Save(pldge, pledgeRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Validator instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
