package sdk

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

var cdc = amino.NewCodec()

func init() {
	registerAmino(cdc)
}

func registerAmino(codec *amino.Codec) {
	// Dependencies:
	crypto.Register(codec)
}
