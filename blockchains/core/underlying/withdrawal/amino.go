package withdrawal

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
)

const (

	// XMNSuiteApplicationsXMNWithdrawal represents the xmnsuite xmn Withdrawal resource
	XMNSuiteApplicationsXMNWithdrawal = "xmnsuite/xmn/Withdrawal"

	// XMNSuiteApplicationsXMNNormalizedWithdrawal represents the xmnsuite xmn Normalized Withdrawal resource
	XMNSuiteApplicationsXMNNormalizedWithdrawal = "xmnsuite/xmn/NormalizedWithdrawal"
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

	// Withdrawal
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Withdrawal)(nil), nil)
		codec.RegisterConcrete(&withdrawal{}, XMNSuiteApplicationsXMNWithdrawal, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedWithdrawal{}, XMNSuiteApplicationsXMNNormalizedWithdrawal, nil)
	}()
}
