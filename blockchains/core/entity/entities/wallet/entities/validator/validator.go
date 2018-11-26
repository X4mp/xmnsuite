package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
)

type validator struct {
	UUID  *uuid.UUID     `json:"id"`
	PKey  tcrypto.PubKey `json:"pubkey"`
	Pldge pledge.Pledge  `json:"pledge"`
}

func createValidator(id *uuid.UUID, pkey tcrypto.PubKey, pldge pledge.Pledge) Validator {
	out := validator{
		UUID:  id,
		PKey:  pkey,
		Pldge: pldge,
	}

	return &out
}

func createValidatorFromNormalized(normalized *normalizedValidator) (Validator, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	pubkey, pubKeyErr := fromEncodedStringToPubKey(normalized.PubKey)
	if pubKeyErr != nil {
		str := fmt.Sprintf("the normalized pubKey (%s) is invalid: %s", normalized.PubKey, pubKeyErr.Error())
		return nil, errors.New(str)
	}

	pldgeIns, pldgeInsErr := pledge.SDKFunc.CreateMetaData().Denormalize()(normalized.Pledge)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	if pldge, ok := pldgeIns.(pledge.Pledge); ok {
		out := createValidator(&id, pubkey, pldge)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", pldgeIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *validator) ID() *uuid.UUID {
	return obj.UUID
}

// PubKey returns the tendermint pubkey
func (obj *validator) PubKey() tcrypto.PubKey {
	return obj.PKey
}

// Pledge returns the pledge
func (obj *validator) Pledge() pledge.Pledge {
	return obj.Pldge
}
