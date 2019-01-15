package project

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
)

const (
	xmnProject           = "xmnsuite/xmn/CommunityProject"
	xmnNormalizedProject = "xmnsuite/xmn/Normalized/CommunityProject"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	wallet.Register(codec)
	proposal.Register(codec)

	// Project
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Project)(nil), nil)
		codec.RegisterConcrete(&project{}, xmnProject, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedProject{}, xmnNormalizedProject, nil)
	}()
}
