package link

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/node"
)

const (

	// XMNSuiteApplicationsXMNLink represents the xmnsuite xmn Link resource
	XMNSuiteApplicationsXMNLink = "xmnsuite/xmn/Link"

	// XMNSuiteApplicationsXMNNormalizedLink represents the xmnsuite xmn Normalized Link resource
	XMNSuiteApplicationsXMNNormalizedLink = "xmnsuite/xmn/NormalizedLink"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	node.Register(codec)

	// Link
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Link)(nil), nil)
		codec.RegisterConcrete(&link{}, XMNSuiteApplicationsXMNLink, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedLink{}, XMNSuiteApplicationsXMNNormalizedLink, nil)
	}()
}
