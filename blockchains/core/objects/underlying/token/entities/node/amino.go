package node

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
)

const (

	// XMNSuiteApplicationsXMNNode represents the xmnsuite xmn Node resource
	XMNSuiteApplicationsXMNNode = "xmnsuite/xmn/Node"

	// XMNSuiteApplicationsXMNNormalizedNode represents the xmnsuite xmn Normalized Node resource
	XMNSuiteApplicationsXMNNormalizedNode = "xmnsuite/xmn/NormalizedNode"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	link.Register(codec)

	// Node
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Node)(nil), nil)
		codec.RegisterConcrete(&node{}, XMNSuiteApplicationsXMNNode, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableNode{}, XMNSuiteApplicationsXMNNormalizedNode, nil)
	}()
}
