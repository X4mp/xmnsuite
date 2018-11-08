package entity

import (
	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	cdc = codec
}
