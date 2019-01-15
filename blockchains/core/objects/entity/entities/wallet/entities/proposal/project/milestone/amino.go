package milestone

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
)

const (
	xmnMilestone           = "xmnsuite/xmn/Milestone"
	xmnNormalizedMilestone = "xmnsuite/xmn/Normalized/Milestone"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	wallet.Register(codec)
	project.Register(codec)
	feature.Register(codec)

	// Milestone
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Milestone)(nil), nil)
		codec.RegisterConcrete(&milestone{}, xmnMilestone, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedMilestone{}, xmnNormalizedMilestone, nil)
	}()
}
