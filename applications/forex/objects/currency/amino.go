package currency

import (
	amino "github.com/tendermint/go-amino"
)

const (
	xmnApplicationsForexCurrency           = "xmn/applications/forex/currency"
	xmnApplicationsForexNormalizedCurrency = "xmn/applications/forex/normalizedCurrency"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Currency
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Currency)(nil), nil)
		codec.RegisterConcrete(&currency{}, xmnApplicationsForexCurrency, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableCurrency{}, xmnApplicationsForexNormalizedCurrency, nil)
	}()
}
