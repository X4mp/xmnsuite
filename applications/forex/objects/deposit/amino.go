package deposit

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

const (
	xmnApplicationsForexDeposit           = "xmn/applications/forex/deposit"
	xmnApplicationsForexNormalizedDeposit = "xmn/applications/forex/normalizedDeposit"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	pledge.Register(codec)
	currency.Register(codec)

	// Bank
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Deposit)(nil), nil)
		codec.RegisterConcrete(&deposit{}, xmnApplicationsForexDeposit, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedDeposit{}, xmnApplicationsForexNormalizedDeposit, nil)
	}()
}
