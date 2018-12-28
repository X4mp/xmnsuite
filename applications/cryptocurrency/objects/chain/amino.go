package chain

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

const (
	xmnCryptocurrencyChain           = "xmn/applications/cryptocurrency/chain"
	xmnCryptocurrencyNormalizedChain = "xmn/applications/cryptocurrency/normalizedChain"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	pledge.Register(codec)

	// Chain
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Chain)(nil), nil)
		codec.RegisterConcrete(&chain{}, xmnCryptocurrencyChain, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedChain{}, xmnCryptocurrencyNormalizedChain, nil)
	}()
}
