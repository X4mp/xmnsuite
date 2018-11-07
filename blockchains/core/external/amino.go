package external

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

const (
	// XMNSuiteApplicationsXMNExternal represents the xmnsuite xmn External resource
	XMNSuiteApplicationsXMNExternal = "xmnsuite/xmn/External"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// External
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*External)(nil), nil)
		codec.RegisterConcrete(&external{}, XMNSuiteApplicationsXMNExternal, nil)
	}()
}
