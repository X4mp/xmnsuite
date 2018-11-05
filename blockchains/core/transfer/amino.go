package transfer

import (
	amino "github.com/tendermint/go-amino"
)

const (

	// XMNSuiteApplicationsXMNTransfer represents the xmnsuite xmn Transfer resource
	XMNSuiteApplicationsXMNTransfer = "xmnsuite/xmn/Transfer"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Token
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Transfer)(nil), nil)
		codec.RegisterConcrete(&transfer{}, XMNSuiteApplicationsXMNTransfer, nil)
	}()
}
