package deposit

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

const (

	// XMNSuiteApplicationsXMNDeposit represents the xmnsuite xmn Deposit resource
	XMNSuiteApplicationsXMNDeposit = "xmnsuite/xmn/Deposit"

	// XMNSuiteApplicationsXMNNormalizedDeposit represents the xmnsuite xmn Normalized Deposit resource
	XMNSuiteApplicationsXMNNormalizedDeposit = "xmnsuite/xmn/Normalized/Deposit"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	token.Register(codec)
	wallet.Register(codec)

	// Deposit
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Deposit)(nil), nil)
		codec.RegisterConcrete(&deposit{}, XMNSuiteApplicationsXMNDeposit, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedDeposit{}, XMNSuiteApplicationsXMNNormalizedDeposit, nil)
	}()
}
