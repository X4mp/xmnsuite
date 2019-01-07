package validator

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Validator represents a validator
type Validator interface {
	ID() *uuid.UUID
	IP() net.IP
	Port() int
	PubKey() tcrypto.PubKey
	Pledge() pledge.Pledge
}

// Normalized represents the normalized validator
type Normalized interface {
}

// Repository represents the validator repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Validator, error)
	RetrieveByPledge(pldge pledge.Pledge) (Validator, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetOrderedByPledgeAmount(index int, amount int) ([]Validator, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	IP     net.IP
	Port   int
	PubKey tcrypto.PubKey
	Pledge pledge.Pledge
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	Store            datastore.DataStore
	EntityRepository entity.Repository
}

// SDKFunc represents the Validator SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Validator
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Validator {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createValidator(params.ID, params.IP, params.Port, params.PubKey, params.Pledge)
		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		met := createMetaData()
		if params.Store != nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.Store)
		}

		out := createRepository(params.EntityRepository, met)
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
				if val, ok := ins.(Validator); ok {
					return []string{
						retrieveAllValidatorsKeyname(),
						retrieveValidatorsByPledgeKeyname(val.Pledge()),
					}, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				// create the metadata and representations:
				metaData := createMetaData()
				pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()

				if val, ok := ins.(Validator); ok {
					// if the validator already exists, return an error:
					_, retValidatorErr := repository.RetrieveByID(metaData, val.ID())
					if retValidatorErr == nil {
						str := fmt.Sprintf("the given Validator (ID: %s) already exists", val.ID().String())
						return errors.New(str)
					}

					// try to retrieve the pledge, return an error if it exists:
					pldge := val.Pledge()
					_, retPledgeErr := repository.RetrieveByID(pledgeRepresentation.MetaData(), pldge.ID())
					if retPledgeErr == nil {
						str := fmt.Sprintf("the validator (ID: %s) contains a pledge (ID: %s) that already exists", val.ID().String(), pldge.ID().String())
						return errors.New(str)
					}

					// save the pledge:
					saveErr := service.Save(pldge, pledgeRepresentation)
					if saveErr != nil {
						return saveErr
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Validator instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
