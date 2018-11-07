package sell

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/framework/user"
)

const (

	// XMNSuiteApplicationsXMNWish represents the xmnsuite xmn Wish resource
	XMNSuiteApplicationsXMNWish = "xmnsuite/xmn/Wish"

	// XMNSuiteApplicationsXMNSell represents the xmnsuite xmn Sell resource
	XMNSuiteApplicationsXMNSell = "xmnsuite/xmn/Sell"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Wish
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Wish)(nil), nil)
		codec.RegisterConcrete(&wish{}, XMNSuiteApplicationsXMNWish, nil)
	}()

	// Sell
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Sell)(nil), nil)
		codec.RegisterConcrete(&sell{}, XMNSuiteApplicationsXMNSell, nil)
	}()
}
