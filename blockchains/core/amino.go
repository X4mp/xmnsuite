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

func Register(codec *amino.Codec) {
	// dependencies:
	account.Register(codec)
	wallet.Register(codec)
	token.Register(codec)
	deposit.Register(codec)
	withdrawal.Register(codec)
	genesis.Register(codec)
	pledge.Register(codec)
	transfer.Register(codec)
	active_request.Register(codec)
	active_vote.Register(codec)
	validator.Register(codec)
	link.Register(codec)
	node.Register(codec)
	developer.Register(codec)
	project.Register(codec)
	milestone.Register(codec)
	task.Register(codec)

	// replace:
	entity.Replace(codec)
	request.Replace(codec)
	active_request.Replace(codec)
	vote.Replace(codec)
	active_vote.Replace(codec)

	cdc = codec
}
