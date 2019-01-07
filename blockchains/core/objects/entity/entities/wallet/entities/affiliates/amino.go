package affiliates

import (
	amino "github.com/tendermint/go-amino"
)

const (
	xmnCloudAffiliate           = "xmn/cloud/affiliate"
	xmnCloudNormalizedAffiliate = "xmn/cloud/normalizedAffiliate"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Affiliate
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Affiliate)(nil), nil)
		codec.RegisterConcrete(&affiliate{}, xmnCloudAffiliate, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedAffiliate{}, xmnCloudNormalizedAffiliate, nil)
	}()
}
