package feature

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

const (
	xmnFeature           = "xmnsuite/xmn/Feature"
	xmnNormalizedFeature = "xmnsuite/xmn/Normalized/Feature"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	project.Register(codec)
	user.Register(codec)

	// Feature
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Feature)(nil), nil)
		codec.RegisterConcrete(&feature{}, xmnFeature, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedFeature{}, xmnNormalizedFeature, nil)
	}()
}
