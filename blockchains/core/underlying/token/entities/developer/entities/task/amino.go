package task

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
)

const (

	// XMNSuiteApplicationsXMNTask represents the xmnsuite xmn Task resource
	XMNSuiteApplicationsXMNTask = "xmnsuite/xmn/Task"

	// XMNSuiteApplicationsXMNNormalizedTask represents the xmnsuite xmn Normalized Task resource
	XMNSuiteApplicationsXMNNormalizedTask = "xmnsuite/xmn/Normalized/Task"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	developer.Register(codec)
	milestone.Register(codec)

	// Task
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Task)(nil), nil)
		codec.RegisterConcrete(&task{}, XMNSuiteApplicationsXMNTask, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedTask{}, XMNSuiteApplicationsXMNNormalizedTask, nil)
	}()
}
