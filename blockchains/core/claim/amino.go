package claim

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/framework/user"
)

const (

	// XMNSuiteApplicationsXMNClaim represents the xmnsuite xmn Claim resource
	XMNSuiteApplicationsXMNClaim = "xmnsuite/xmn/Claim"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Claim
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Claim)(nil), nil)
		codec.RegisterConcrete(&claim{}, XMNSuiteApplicationsXMNClaim, nil)
	}()
}
