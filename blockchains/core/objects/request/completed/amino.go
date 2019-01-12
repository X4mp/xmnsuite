package completed

import (
	amino "github.com/tendermint/go-amino"
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

const (

	// XMNSuiteBlockchainsFrameworkCompletedRequest represents the xmnsuite core completed request
	XMNSuiteBlockchainsFrameworkCompletedRequest = "xmnsuite/blockchains/core/CompletedRequest"

	// XMNSuiteBlockchainsFrameworkNormalizedCompletedRequest represents the xmnsuite core normalized completed request
	XMNSuiteBlockchainsFrameworkNormalizedCompletedRequest = "xmnsuite/blockchains/core/NormalizedCompletedRequest"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies:
	prev_request.Register(codec)

	// Request
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Request)(nil), nil)
		codec.RegisterConcrete(&request{}, XMNSuiteBlockchainsFrameworkCompletedRequest, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedRequest{}, XMNSuiteBlockchainsFrameworkNormalizedCompletedRequest, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
