package task

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

const (
	xmnTask           = "xmnsuite/xmn/Task"
	xmnNormalizedTask = "xmnsuite/xmn/Normalized/Task"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	milestone.Register(codec)
	user.Register(codec)

	// Task
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Task)(nil), nil)
		codec.RegisterConcrete(&task{}, xmnTask, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedTask{}, xmnNormalizedTask, nil)
	}()
}
