package core

import (
	amino "github.com/tendermint/go-amino"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	transfer "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/vote"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	developer "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	milestone "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/milestone"
	project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/project"
	task "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/task"
	link "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	node "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
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
	node.Register(cdc)
	developer.Register(cdc)
	project.Register(cdc)
	milestone.Register(cdc)
	task.Register(cdc)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
	vote.Replace(cdc)
}
