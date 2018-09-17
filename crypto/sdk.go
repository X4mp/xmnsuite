package crypto

import (
	"github.com/dedis/kyber"
)

// PrivateKey represents a private key
type PrivateKey interface {
	PublicKey() PublicKey
	Sign(msg string) Signature
	RingSign(msg string, ringPubKeys []PublicKey) (RingSignature, error)
	String() string
}

// PublicKey represents the public key
type PublicKey interface {
	Point() kyber.Point
	Equals(pubKey PublicKey) bool
	String() string
}

// Signature represents a signature
type Signature interface {
	PublicKey(msg string) PublicKey
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

// CreatePubKeyParams represents the CreatePubKey func params
type CreatePubKeyParams struct {
	PubKeyAsString string
}

// CreateSigParams represents the CreateSig params
type CreateSigParams struct {
	SigAsString string
}

// SDKFunc represents the crypto SDK func
var SDKFunc = struct {
	GenPK        func() PrivateKey
	CreatePK     func(params CreatePKParams) PrivateKey
	CreatePubKey func(params CreatePubKeyParams) PublicKey
	CreateSig    func(params CreateSigParams) Signature
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

	CreatePubKey: func(params CreatePubKeyParams) PublicKey {
		pubKey, pubKeyErr := createPublicKeyFromString(params.PubKeyAsString)
		if pubKeyErr != nil {
			panic(pubKeyErr)
		}

		return pubKey
	},

	CreateSig: func(params CreateSigParams) Signature {
		sig, sigErr := createSignatureFromString(params.SigAsString)
		if sigErr != nil {
			panic(sigErr)
		}

		return sig
	},
}
