package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

// Validator represents a validator
type Validator interface {
	ID() *uuid.UUID
	PubKey() tcrypto.PubKey
	Pledge() pledge.Pledge
}

// SDKFunc represents the Validator SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
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
		})
	},
}
