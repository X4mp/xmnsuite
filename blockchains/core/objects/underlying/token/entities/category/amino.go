package category

import (
	amino "github.com/tendermint/go-amino"
)

const (
	xmnCategory           = "xmnsuite/xmn/Category"
	xmnNormalizedCategory = "xmnsuite/xmn/Normalized/Category"
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
		codec.RegisterConcrete(&category{}, xmnCategory, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedCategory{}, xmnNormalizedCategory, nil)
	}()
}
