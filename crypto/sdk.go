package crypto

import (
	"github.com/dedis/kyber"
)

// PrivateKey represents a private key
type PrivateKey interface {
	PublicKey() kyber.Point
	Sign(msg string) Signature
	RingSign(msg string, ringPubKeys []kyber.Point) (RingSignature, error)
	String() string
}

// Signature represents a signature
type Signature interface {
	PublicKey(msg string) kyber.Point
	Verify(msg string) bool
	String() string
}

// RingSignature represents a RingSignature
type RingSignature interface {
	Verify(msg string) bool
	String() string
}

// CreatePKParams represents the CreatePK func params
type CreatePKParams struct {
	PKAsString string
}

// SDKFunc represents the crypto SDK func
var SDKFunc = struct {
	GenPK    func() PrivateKey
	CreatePK func(params CreatePKParams) PrivateKey
}{
	GenPK: func() PrivateKey {
		return createPrivateKey()
	},

	CreatePK: func(params CreatePKParams) PrivateKey {
		if params.PKAsString == "" {
			return createPrivateKey()
		}

		pk, pkErr := createPrivateKeyFromString(params.PKAsString)
		if pkErr != nil {
			panic(pkErr)
		}

		return pk
	},
}
