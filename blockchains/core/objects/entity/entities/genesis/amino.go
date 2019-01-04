package genesis

import (
	amino "github.com/tendermint/go-amino"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	user "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

const (

	// XMNSuiteApplicationsXMNGenesis represents the xmnsuite xmn Genesis resource
	XMNSuiteApplicationsXMNGenesis = "xmnsuite/xmn/Genesis"

	// XMNSuiteApplicationsXMNNormalizedGenesis represents the xmnsuite xmn Normalized Genesis resource
	XMNSuiteApplicationsXMNNormalizedGenesis = "xmnsuite/xmn/Normalized/Genesis"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	deposit.Register(codec)
	user.Register(codec)

	// Genesis
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Genesis)(nil), nil)
		codec.RegisterConcrete(&genesis{}, XMNSuiteApplicationsXMNGenesis, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedGenesis{}, XMNSuiteApplicationsXMNNormalizedGenesis, nil)
	}()
}
