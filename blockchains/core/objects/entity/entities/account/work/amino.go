package work

import (
	amino "github.com/tendermint/go-amino"
	user "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
)

const (

	// XMNSuiteApplicationsAccountWork represents the xmnsuite xmn Account Work resource
	XMNSuiteApplicationsAccountWork = "xmnsuite/xmn/Account/Work"

	// XMNSuiteApplicationsAccountNormalizedWork represents the xmnsuite xmn Normalized Account Work resource
	XMNSuiteApplicationsAccountNormalizedWork = "xmnsuite/xmn/Account/Normalized/Work"
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

	// Work
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Work)(nil), nil)
		codec.RegisterConcrete(&work{}, XMNSuiteApplicationsAccountWork, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedWork{}, XMNSuiteApplicationsAccountNormalizedWork, nil)
	}()
}
