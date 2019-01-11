package pledge

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	mils_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
)

const (
	xmnTask           = "xmnsuite/xmn/PledgeTask"
	xmnNormalizedTask = "xmnsuite/xmn/Normalized/PledgeTask"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	pledge.Register(codec)
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
