package category

import (
	amino "github.com/tendermint/go-amino"
)

const (
	xmnApplicationsForexCategory           = "xmn/applications/forex/category"
	xmnApplicationsForexNormalizedCategory = "xmn/applications/forex/normalizedCategory"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Category
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Category)(nil), nil)
		codec.RegisterConcrete(&category{}, xmnApplicationsForexCategory, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableCategory{}, xmnApplicationsForexNormalizedCategory, nil)
	}()
}
