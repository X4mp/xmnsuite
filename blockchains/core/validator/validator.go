package validator

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/pledge"
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
