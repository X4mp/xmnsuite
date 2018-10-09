package tendermint

import (
	amino "github.com/tendermint/go-amino"
	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	//Dependencies:
	func() {
		applications.Register(codec)
		crypto.Register(codec)
	}()

	// PublicKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*tcrypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.PubKeyAminoRoute, nil)
	}()

	// PrivateKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*tcrypto.PrivKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.PrivKeyAminoRoute, nil)
	}()
}
