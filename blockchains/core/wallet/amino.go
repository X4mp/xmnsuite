package wallet

import (
	amino "github.com/tendermint/go-amino"
	applications "github.com/xmnservices/xmnsuite/routers"
)

const (

	// XMNSuiteApplicationsXMNWallet represents the xmnsuite xmn Wallet resource
	XMNSuiteApplicationsXMNWallet = "xmnsuite/xmn/Wallet"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	applications.Register(codec)

	// Wallet
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Wallet)(nil), nil)
		codec.RegisterConcrete(&wallet{}, XMNSuiteApplicationsXMNWallet, nil)
	}()
}
