package vote

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
)

const (

	// XMNSuiteBlockchainsFrameworkVote represents the xmnsuite core Vote
	XMNSuiteBlockchainsFrameworkVote = "xmnsuite/blockchains/core/Vote"

	// XMNSuiteBlockchainsFrameworkNormalizedVote represents the xmnsuite core NormalizedVote
	XMNSuiteBlockchainsFrameworkNormalizedVote = "xmnsuite/blockchains/core/NormalizedVote"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)
	request.Register(codec)

	// Vote
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Vote)(nil), nil)
		codec.RegisterConcrete(&vote{}, XMNSuiteBlockchainsFrameworkVote, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedVote{}, XMNSuiteBlockchainsFrameworkNormalizedVote, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
