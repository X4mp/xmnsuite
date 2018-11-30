package applications

import (
	"net"

	"github.com/tendermint/tendermint/crypto"
)

type validator struct {
	ipAddress net.IP
	pubKey    crypto.PubKey
	pow       int64
}

func createValidator(ipAddress net.IP, pubKey crypto.PubKey, pow int64) Validator {
	out := validator{
		ipAddress: ipAddress,
		pubKey:    pubKey,
		pow:       pow,
	}

	return &out
}

// IP returns the ip address
func (obj *validator) IP() net.IP {
	return obj.ipAddress
}

// PubKey returns the pubKey
func (obj *validator) PubKey() crypto.PubKey {
	return obj.pubKey
}

// Power returns the power
func (obj *validator) Power() int64 {
	return obj.pow
}
