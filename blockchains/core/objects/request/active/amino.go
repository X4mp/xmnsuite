package active

import (
	amino "github.com/tendermint/go-amino"
	prev_req "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

const (

	// XMNSuiteBlockchainsCoreRequestActive represents the xmnsuite core Request
	XMNSuiteBlockchainsCoreRequestActive = "xmnsuite/blockchains/core/Request/ActiveRequest"

	// XMNSuiteBlockchainsCoreNormalizedRequestActive represents the xmnsuite core NormalizedRequest
	XMNSuiteBlockchainsCoreNormalizedRequestActive = "xmnsuite/blockchains/core/Request/NormalizedActiveRequest"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	prev_req.Register(codec)

	// Request
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Request)(nil), nil)
		codec.RegisterConcrete(&request{}, XMNSuiteBlockchainsCoreRequestActive, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedRequest{}, XMNSuiteBlockchainsCoreNormalizedRequestActive, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
