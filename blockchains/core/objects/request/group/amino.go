package group

import (
	amino "github.com/tendermint/go-amino"
)

const (

	// XMNSuiteBlockchainsFrameworkGroup represents the xmnsuite core Group
	XMNSuiteBlockchainsFrameworkGroup = "xmnsuite/blockchains/core/Group"

	// XMNSuiteBlockchainsFrameworkNormalizedGroup represents the xmnsuite core NormalizedVote
	XMNSuiteBlockchainsFrameworkNormalizedGroup = "xmnsuite/blockchains/core/NormalizedGroup"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Group
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Group)(nil), nil)
		codec.RegisterConcrete(&group{}, XMNSuiteBlockchainsFrameworkGroup, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableGroup{}, XMNSuiteBlockchainsFrameworkNormalizedGroup, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
