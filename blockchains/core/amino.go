package core

import (
	amino "github.com/tendermint/go-amino"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/genesis"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/pledge"
	request "github.com/xmnservices/xmnsuite/blockchains/core/request"
	token "github.com/xmnservices/xmnsuite/blockchains/core/token"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

var cdc = amino.NewCodec()

func init() {
	// dependencies:
	wallet.Register(cdc)
	token.Register(cdc)
	deposit.Register(cdc)
	genesis.Register(cdc)
	pledge.Register(cdc)
	request.Register(cdc)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
}
