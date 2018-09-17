package crypto

import (
	amino "github.com/tendermint/go-amino"
)

const (
	// XMNSuiteCryptoPrivateKey represents the xmnsuite crypto PrivateKey
	XMNSuiteCryptoPrivateKey = "xmnsuite/crypto/PrivateKey"

	// XMNSuiteCryptoPublicKey represents the xmnsuite crypto PublicKey
	XMNSuiteCryptoPublicKey = "xmnsuite/crypto/PublicKey"

	// XMNSuiteCryptoSignature represents the xmnsuite crypto Signature
	XMNSuiteCryptoSignature = "xmnsuite/crypto/Signature"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// PrivateKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*PrivateKey)(nil), nil)
		codec.RegisterConcrete(&privateKey{}, XMNSuiteCryptoPrivateKey, nil)
	}()

	// PublicKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*PublicKey)(nil), nil)
		codec.RegisterConcrete(&pubKey{}, XMNSuiteCryptoPublicKey, nil)
	}()

	// Signature
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Signature)(nil), nil)
		codec.RegisterConcrete(&signature{}, XMNSuiteCryptoSignature, nil)
	}()
}
