package fiatchain

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
)

const (
	xmnApplicationsForexFiatChain           = "xmn/applications/forex/fiatChain"
	xmnApplicationsForexNormalizedFiatChain = "xmn/applications/forex/normalizedFiatChain"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	deposit.Register(codec)

	// Bank
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*FiatChain)(nil), nil)
		codec.RegisterConcrete(&fiatChain{}, xmnApplicationsForexFiatChain, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedFiatChain{}, xmnApplicationsForexNormalizedFiatChain, nil)
	}()
}
