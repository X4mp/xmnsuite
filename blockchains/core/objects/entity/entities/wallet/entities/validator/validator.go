package validator

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
)

type validator struct {
	UUID      *uuid.UUID     `json:"id"`
	IPAddress net.IP         `json:"ip"`
	Prt       int            `json:"port"`
	PKey      tcrypto.PubKey `json:"pubkey"`
	Pldge     pledge.Pledge  `json:"pledge"`
}

func createValidator(id *uuid.UUID, ipAddress net.IP, port int, pkey tcrypto.PubKey, pldge pledge.Pledge) Validator {
	out := validator{
		UUID:      id,
		IPAddress: ipAddress,
		Prt:       port,
		PKey:      pkey,
		Pldge:     pldge,
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
		ip := net.ParseIP(normalized.IP)
		out := createValidator(&id, ip, normalized.Port, pubkey, pldge)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Validator instance", pldgeIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *validator) ID() *uuid.UUID {
	return obj.UUID
}

// IP returns the IP
func (obj *validator) IP() net.IP {
	return obj.IPAddress
}

// Port returns the port
func (obj *validator) Port() int {
	return obj.Prt
}

// PubKey returns the tendermint pubkey
func (obj *validator) PubKey() tcrypto.PubKey {
	return obj.PKey
}

// Pledge returns the pledge
func (obj *validator) Pledge() pledge.Pledge {
	return obj.Pldge
}
