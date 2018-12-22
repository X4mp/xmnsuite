package keyname

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
)

const (

	// XMNSuiteBlockchainsFrameworkKeyname represents the xmnsuite core Keyname
	XMNSuiteBlockchainsFrameworkKeyname = "xmnsuite/blockchains/core/Keyname"

	// XMNSuiteBlockchainsFrameworkNormalizedKeyname represents the xmnsuite core NormalizedVote
	XMNSuiteBlockchainsFrameworkNormalizedKeyname = "xmnsuite/blockchains/core/NormalizedKeyname"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	group.Register(codec)

	// Keyname
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Keyname)(nil), nil)
		codec.RegisterConcrete(&keyname{}, XMNSuiteBlockchainsFrameworkKeyname, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedKeyname{}, XMNSuiteBlockchainsFrameworkNormalizedKeyname, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
