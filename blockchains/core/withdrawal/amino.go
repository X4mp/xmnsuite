package withdrawal

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
)

const (

	// XMNSuiteApplicationsXMNWithdrawal represents the xmnsuite xmn Withdrawal resource
	XMNSuiteApplicationsXMNWithdrawal = "xmnsuite/xmn/Withdrawal"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Withdrawal
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Withdrawal)(nil), nil)
		codec.RegisterConcrete(&withdrawal{}, XMNSuiteApplicationsXMNWithdrawal, nil)
	}()
}
