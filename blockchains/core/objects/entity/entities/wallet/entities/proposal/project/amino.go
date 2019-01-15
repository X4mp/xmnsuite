package project

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
)

const (
	xmnProject           = "xmnsuite/xmn/Project"
	xmnNormalizedProject = "xmnsuite/xmn/Normalized/Project"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	wallet.Register(codec)
	approved_project.Register(codec)

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
