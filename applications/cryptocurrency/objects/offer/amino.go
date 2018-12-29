package offer

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

const (
	xmnCryptocurrencyOffer           = "xmn/applications/cryptocurrency/offer"
	xmnCryptocurrencyNormalizedOffer = "xmn/applications/cryptocurrency/normalizedOffer"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	pledge.Register(codec)
	address.Register(codec)

	// Offer
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Offer)(nil), nil)
		codec.RegisterConcrete(&offer{}, xmnCryptocurrencyOffer, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedOffer{}, xmnCryptocurrencyNormalizedOffer, nil)
	}()
}
