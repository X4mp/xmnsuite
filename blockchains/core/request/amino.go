package request

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

const (

	// XMNSuiteBlockchainsCoreRequest represents the xmnsuite core Request
	XMNSuiteBlockchainsCoreRequest = "xmnsuite/blockchains/core/Request"

	// XMNSuiteBlockchainsCoreNormalizedRequest represents the xmnsuite core NormalizedRequest
	XMNSuiteBlockchainsCoreNormalizedRequest = "xmnsuite/blockchains/core/NormalizedRequest"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Request
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Request)(nil), nil)
		codec.RegisterConcrete(&request{}, XMNSuiteBlockchainsCoreRequest, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedRequest{}, XMNSuiteBlockchainsCoreNormalizedRequest, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
