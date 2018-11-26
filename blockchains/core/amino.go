package core

import (
	amino "github.com/tendermint/go-amino"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/validator"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	request "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/pledge"
	transfer "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/transfer"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/vote"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	link "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/request/entities/link"
)

var cdc = amino.NewCodec()

func init() {
	// dependencies:
	wallet.Register(cdc)
	token.Register(cdc)
	deposit.Register(cdc)
	genesis.Register(cdc)
	pledge.Register(cdc)
	transfer.Register(cdc)
	request.Register(cdc)
	vote.Register(cdc)
	validator.Register(cdc)
	link.Register(cdc)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
	vote.Replace(cdc)
}
