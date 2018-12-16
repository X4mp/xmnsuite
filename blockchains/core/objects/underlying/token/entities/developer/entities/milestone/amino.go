package milestone

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/project"
)

const (

	// XMNSuiteApplicationsXMNMilestone represents the xmnsuite xmn Milestone resource
	XMNSuiteApplicationsXMNMilestone = "xmnsuite/xmn/Milestone"

	// XMNSuiteApplicationsXMNNormalizedMilestone represents the xmnsuite xmn Normalized Milestone resource
	XMNSuiteApplicationsXMNNormalizedMilestone = "xmnsuite/xmn/Normalized/Milestone"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	project.Register(codec)

	// Milestone
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Milestone)(nil), nil)
		codec.RegisterConcrete(&milestone{}, XMNSuiteApplicationsXMNMilestone, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedMilestone{}, XMNSuiteApplicationsXMNNormalizedMilestone, nil)
	}()
}
