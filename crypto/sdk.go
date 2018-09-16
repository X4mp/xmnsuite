package crypto

import (
	"github.com/dedis/kyber"
)

// PrivateKey represents a private key
type PrivateKey interface {
	PublicKey() kyber.Point
	Sign(msg string) Signature
	RingSign(msg string, ringPubKeys []kyber.Point) RingSignature
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
}

// CreatePKParams represents the CreatePK func params
type CreatePKParams struct {
	PKAsString string
}

// CreateSigParams represents the CreateSig func params
type CreateSigParams struct {
	SigAsString string
}

// SDKFunc represents the crypto SDK func
var SDKFunc = struct {
	GenPK     func() PrivateKey
	CreatePK  func(params CreatePKParams) PrivateKey
	CreateSig func(params CreateSigParams) Signature
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
	CreateSig: func(params CreateSigParams) Signature {
		sig, sigErr := createSignatureFromString(params.SigAsString)
		if sigErr != nil {
			panic(sigErr)
		}

		return sig
	},
}
