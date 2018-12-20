package web

import (
	amino "github.com/tendermint/go-amino"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
)

var cdc = amino.NewCodec()

func init() {
	// dependencies:
	core.Register(cdc)
}
