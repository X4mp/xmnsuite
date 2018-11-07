package pledge

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

const (

	// XMNSuiteApplicationsXMNPledge represents the xmnsuite xmn Pledge resource
	XMNSuiteApplicationsXMNPledge = "xmnsuite/xmn/Pledge"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Pledge
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Pledge)(nil), nil)
		codec.RegisterConcrete(&pledge{}, XMNSuiteApplicationsXMNPledge, nil)
	}()
}
