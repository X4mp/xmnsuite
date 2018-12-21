package web

import (
	amino "github.com/tendermint/go-amino"
	bank "github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	category "github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	currency "github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	deposit "github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
	fiatchain "github.com/xmnservices/xmnsuite/applications/forex/objects/fiatchain"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
)

var cdc = amino.NewCodec()

func init() {
	// dependencies:
	core.Register(cdc)
	category.Register(cdc)
	bank.Register(cdc)
	currency.Register(cdc)
	deposit.Register(cdc)
	fiatchain.Register(cdc)
}
