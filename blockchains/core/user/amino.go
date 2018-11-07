package user

import (
	amino "github.com/tendermint/go-amino"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	applications "github.com/xmnservices/xmnsuite/routers"
)

const (

	// XMNSuiteApplicationsXMNUser represents the xmnsuite xmn User resource
	XMNSuiteApplicationsXMNUser = "xmnsuite/xmn/User"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	applications.Register(codec)
	wallet.Register(codec)

	// User
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*User)(nil), nil)
		codec.RegisterConcrete(&user{}, XMNSuiteApplicationsXMNUser, nil)
	}()
}
