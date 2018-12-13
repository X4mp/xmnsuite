package main

import (
	amino "github.com/tendermint/go-amino"
	forex "github.com/xmnservices/xmnsuite/applications/forex"
)

var cdc = amino.NewCodec()

func init() {
	registerAmino(cdc)
}

func registerAmino(codec *amino.Codec) {
	forex.Register(codec)
}
