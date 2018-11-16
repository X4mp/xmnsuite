package link

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet/request/entities/user"
)

const (

	// XMNSuiteApplicationsXMNLink represents the xmnsuite xmn Link resource
	XMNSuiteApplicationsXMNLink = "xmnsuite/xmn/Link"

	// XMNSuiteApplicationsXMNNode represents the xmnsuite xmn Node resource
	XMNSuiteApplicationsXMNNode = "xmnsuite/xmn/Node"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Link
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Link)(nil), nil)
		codec.RegisterConcrete(&link{}, XMNSuiteApplicationsXMNLink, nil)
	}()

	// Deposit
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Node)(nil), nil)
		codec.RegisterConcrete(&node{}, XMNSuiteApplicationsXMNNode, nil)
	}()
}
