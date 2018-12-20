package account

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	user "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
)

const (

	// XMNSuiteApplicationsAccount represents the xmnsuite xmn Account resource
	XMNSuiteApplicationsAccount = "xmnsuite/xmn/Account"

	// XMNSuiteApplicationsNormalizedAccount represents the xmnsuite xmn Normalized Account resource
	XMNSuiteApplicationsNormalizedAccount = "xmnsuite/xmn/Normalized/Account"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	work.Register(codec)
	user.Register(codec)

	// Work
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Account)(nil), nil)
		codec.RegisterConcrete(&account{}, XMNSuiteApplicationsAccount, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedAccount{}, XMNSuiteApplicationsNormalizedAccount, nil)
	}()
}
