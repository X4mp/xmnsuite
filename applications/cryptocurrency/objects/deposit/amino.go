package deposit

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
)

const (
	xmnCryptocurrencyDeposit           = "xmn/applications/cryptocurrency/deposit"
	xmnCryptocurrencyNormalizedDeposit = "xmn/applications/cryptocurrency/normalizedDeposit"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	offer.Register(codec)
	address.Register(codec)

	// Deposit
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Deposit)(nil), nil)
		codec.RegisterConcrete(&deposit{}, xmnCryptocurrencyDeposit, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedDeposit{}, xmnCryptocurrencyNormalizedDeposit, nil)
	}()
}
