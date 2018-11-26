package validator

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
)

const (

	// XMNSuiteApplicationsXMNValidator represents the xmnsuite xmn Validator resource
	XMNSuiteApplicationsXMNValidator = "xmnsuite/xmn/Validator"

	// XMNSuiteApplicationsXMNNormalizedValidator represents the xmnsuite xmn NormalizedValidator resource
	XMNSuiteApplicationsXMNNormalizedValidator = "xmnsuite/xmn/NormalizedValidator"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	pledge.Register(codec)

	// crypto.PubKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.PubKeyAminoRoute, nil)
	}()

	// Validator
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Validator)(nil), nil)
		codec.RegisterConcrete(&validator{}, XMNSuiteApplicationsXMNValidator, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedValidator{}, XMNSuiteApplicationsXMNNormalizedValidator, nil)
	}()
}
