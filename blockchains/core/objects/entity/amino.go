package entity

import (
	amino "github.com/tendermint/go-amino"
)

const (
	aminoEntityPartialSet = "xmnsuite/xmn/EntityPartialSet"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	cdc = codec
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedPartialSet{}, aminoEntityPartialSet, nil)
	}()
}
