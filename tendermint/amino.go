package tendermint

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// PublicKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)
	}()

	// PrivateKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PrivKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.Ed25519PrivKeyAminoRoute, nil)
	}()
}
