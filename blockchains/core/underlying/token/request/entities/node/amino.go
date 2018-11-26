package node

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
)

const (

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

	// Node
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Node)(nil), nil)
		codec.RegisterConcrete(&node{}, XMNSuiteApplicationsXMNNode, nil)
	}()
}
