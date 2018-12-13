package bank

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
)

const (
	xmnApplicationsForexBank           = "xmn/applications/forex/bank"
	xmnApplicationsForexNormalizedBank = "xmn/applications/forex/normalizedBank"
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
		codec.RegisterInterface((*Bank)(nil), nil)
		codec.RegisterConcrete(&bank{}, xmnApplicationsForexBank, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedBank{}, xmnApplicationsForexNormalizedBank, nil)
	}()
}
