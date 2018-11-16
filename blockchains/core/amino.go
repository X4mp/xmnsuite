package core

import (
	amino "github.com/tendermint/go-amino"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/pledge"
	request "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	token "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/validator"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/vote"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
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
	vote.Register(cdc)
	validator.Register(cdc)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
	vote.Replace(cdc)
}
