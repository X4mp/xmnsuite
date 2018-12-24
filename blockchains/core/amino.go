package core

import (
	amino "github.com/tendermint/go-amino"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	account "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	transfer "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/transfer"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/validator"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	developer "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer"
	milestone "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/milestone"
	project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/project"
	task "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/developer/entities/task"
	link "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	node "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	withdrawal "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

func Register(cdc *amino.Codec) {
	// dependencies:
	account.Register(cdc)
	wallet.Register(cdc)
	token.Register(cdc)
	deposit.Register(cdc)
	withdrawal.Register(cdc)
	genesis.Register(cdc)
	pledge.Register(cdc)
	transfer.Register(cdc)
	active_request.Register(cdc)
	active_vote.Register(cdc)
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
	active_request.Replace(cdc)
	vote.Replace(cdc)
	active_vote.Replace(cdc)
}
