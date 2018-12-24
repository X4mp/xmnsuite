package active

import (
	amino "github.com/tendermint/go-amino"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
)

const (

	// XMNSuiteBlockchainsCoreRequestActiveVote represents the xmnsuite core Request
	XMNSuiteBlockchainsCoreRequestActiveVote = "xmnsuite/blockchains/core/Request/ActiveRequest/Vote"

	// XMNSuiteBlockchainsCoreRequestActiveNormalizedVote represents the xmnsuite core NormalizedRequest
	XMNSuiteBlockchainsCoreRequestActiveNormalizedVote = "xmnsuite/blockchains/core/Request/ActiveRequest/NormalizedVote"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	core_vote.Register(codec)

	// vote
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Vote)(nil), nil)
		codec.RegisterConcrete(&vote{}, XMNSuiteBlockchainsCoreRequestActiveVote, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedVote{}, XMNSuiteBlockchainsCoreRequestActiveNormalizedVote, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
