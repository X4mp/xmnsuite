package deposit

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

const (

	// XMNSuiteApplicationsXMNDeposit represents the xmnsuite xmn Deposit resource
	XMNSuiteApplicationsXMNDeposit = "xmnsuite/xmn/Deposit"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Deposit
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Deposit)(nil), nil)
		codec.RegisterConcrete(&deposit{}, XMNSuiteApplicationsXMNDeposit, nil)
	}()
}
