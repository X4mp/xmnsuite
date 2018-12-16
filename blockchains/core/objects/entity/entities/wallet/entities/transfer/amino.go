package transfer

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

const (

	// XMNSuiteApplicationsXMNTransfer represents the xmnsuite xmn Transfer resource
	XMNSuiteApplicationsXMNTransfer = "xmnsuite/xmn/Transfer"

	// XMNSuiteApplicationsXMNNornalizedTransfer represents the xmnsuite xmn Normalized Transfer resource
	XMNSuiteApplicationsXMNNornalizedTransfer = "xmnsuite/xmn/NormalizedTransfer"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	withdrawal.Register(codec)
	deposit.Register(codec)

	// Transfer
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Transfer)(nil), nil)
		codec.RegisterConcrete(&transfer{}, XMNSuiteApplicationsXMNTransfer, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedTransfer{}, XMNSuiteApplicationsXMNNornalizedTransfer, nil)
	}()
}
