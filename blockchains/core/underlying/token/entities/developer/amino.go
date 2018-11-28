package developer

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

const (

	// XMNSuiteApplicationsXMNDeveloper represents the xmnsuite xmn Developer resource
	XMNSuiteApplicationsXMNDeveloper = "xmnsuite/xmn/Developer"

	// XMNSuiteApplicationsXMNNormalizedDeveloper represents the xmnsuite xmn Normalized Developer resource
	XMNSuiteApplicationsXMNNormalizedDeveloper = "xmnsuite/xmn/Normalized/Developer"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	user.Register(codec)
	pledge.Register(codec)

	// Developer
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Developer)(nil), nil)
		codec.RegisterConcrete(&developer{}, XMNSuiteApplicationsXMNDeveloper, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedDeveloper{}, XMNSuiteApplicationsXMNNormalizedDeveloper, nil)
	}()
}
