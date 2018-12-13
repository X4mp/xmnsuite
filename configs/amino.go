package configs

import (
	amino "github.com/tendermint/go-amino"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (

	// xmnSuiteConfigs represents the xmnsuite configs resource
	xmnSuiteConfigs = "xmnsuite/xmn/configs"

	// xmnSuiteNormalizedConfigs represents the xmnsuite normalized configs resource
	xmnSuiteNormalizedConfigs = "xmnsuite/xmn/configs/normalized"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// crypto.PrivKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*tcrypto.PrivKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.PrivKeyAminoRoute, nil)
	}()

	// Validator
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Configs)(nil), nil)
		codec.RegisterConcrete(&configs{}, xmnSuiteConfigs, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableConfigs{}, xmnSuiteNormalizedConfigs, nil)
	}()
}
