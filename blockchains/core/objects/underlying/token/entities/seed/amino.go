package seed

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
)

const (

	// XMNSuiteApplicationsXMNSeed represents the xmnsuite xmn Seed resource
	XMNSuiteApplicationsXMNSeed = "xmnsuite/xmn/Seed"

	// XMNSuiteApplicationsXMNNormalizedSeed represents the xmnsuite xmn Normalized Seed resource
	XMNSuiteApplicationsXMNNormalizedSeed = "xmnsuite/xmn/NormalizedSeed"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	link.Register(codec)

	// Seed
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Seed)(nil), nil)
		codec.RegisterConcrete(&seed{}, XMNSuiteApplicationsXMNSeed, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableSeed{}, XMNSuiteApplicationsXMNNormalizedSeed, nil)
	}()
}
