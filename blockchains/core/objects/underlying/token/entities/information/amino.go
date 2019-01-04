package information

import (
	amino "github.com/tendermint/go-amino"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	user "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

const (

	// XMNSuiteApplicationsXMNInformation represents the xmnsuite xmn Information resource
	XMNSuiteApplicationsXMNInformation = "xmnsuite/xmn/Information"

	// XMNSuiteApplicationsXMNNormalizedInformation represents the xmnsuite xmn Normalized Information resource
	XMNSuiteApplicationsXMNNormalizedInformation = "xmnsuite/xmn/Normalized/Information"
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

	// Information
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Information)(nil), nil)
		codec.RegisterConcrete(&information{}, XMNSuiteApplicationsXMNInformation, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedInformation{}, XMNSuiteApplicationsXMNNormalizedInformation, nil)
	}()
}
