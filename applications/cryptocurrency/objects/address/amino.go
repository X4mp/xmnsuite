package address

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

const (
	xmnCryptocurrencyAddress           = "xmn/applications/cryptocurrency/address"
	xmnCryptocurrencyNormalizedAddress = "xmn/applications/cryptocurrency/normalizedAddress"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	pledge.Register(codec)

	// Address
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Address)(nil), nil)
		codec.RegisterConcrete(&address{}, xmnCryptocurrencyAddress, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedAddress{}, xmnCryptocurrencyNormalizedAddress, nil)
	}()
}
