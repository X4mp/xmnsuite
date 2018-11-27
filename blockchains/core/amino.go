package core

import (
	amino "github.com/tendermint/go-amino"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	transfer "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/transfer"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/validator"
	request "github.com/xmnservices/xmnsuite/blockchains/core/request"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/request/vote"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	developer "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer"
	milestone "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/milestone"
	project "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/developer/entities/project"
	link "github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
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
	developer.Register(cdc)
	project.Register(cdc)
	milestone.Register(cdc)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
	vote.Replace(cdc)
}
