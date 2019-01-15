package proposal

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

const (
	xmnProposal           = "xmnsuite/xmn/Proposal"
	xmnNormalizedProposal = "xmnsuite/xmn/Normalized/Proposal"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	category.Register(codec)

	// Proposal
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Proposal)(nil), nil)
		codec.RegisterConcrete(&proposal{}, xmnProposal, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedProposal{}, xmnNormalizedProposal, nil)
	}()
}
