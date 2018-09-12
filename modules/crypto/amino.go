package crypto

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var cdc = amino.NewCodec()

func init() {
	registerAmino(cdc)
}

func registerAmino(codec *amino.Codec) {
	// crypto.PubKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)
	}()

	// crypto.PrivKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PrivKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.Ed25519PrivKeyAminoRoute, nil)
	}()
}
