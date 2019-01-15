package completed

import (
	amino "github.com/tendermint/go-amino"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

const (
	xmnTask           = "xmnsuite/xmn/CompletedTask"
	xmnNormalizedTask = "xmnsuite/xmn/Normalized/CompletedTask"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	mils_task.Register(codec)

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
