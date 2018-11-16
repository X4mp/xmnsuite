package applications

import (
	"github.com/tendermint/tendermint/crypto"
)

type validator struct {
	pubKey crypto.PubKey
	pow    int64
}

func createValidator(pubKey crypto.PubKey, pow int64) Validator {
	out := validator{
		pubKey: pubKey,
		pow:    pow,
	}

	return &out
}

// PubKey returns the pubKey
func (obj *validator) PubKey() crypto.PubKey {
	return obj.pubKey
}

// Power returns the power
func (obj *validator) Power() int64 {
	return obj.pow
}
