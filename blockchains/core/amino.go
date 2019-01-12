package core

import (
	amino "github.com/tendermint/go-amino"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	genesis "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	wallet "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	affiliates "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	pledge "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	project "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	feature "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/feature"
	milestone "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone"
	task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task"
	completed_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task/completed"
	pledge_task "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project/milestone/task/pledge"
	transfer "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	user "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	validator "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/completed"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	category "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
	link "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	node "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	coomunity_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
	withdrawal "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register register the amino codec
func Register(codec *amino.Codec) {
	// dependencies:
	wallet.Register(codec)
	user.Register(codec)
	affiliates.Register(codec)
	token.Register(codec)
	deposit.Register(codec)
	withdrawal.Register(codec)
	genesis.Register(codec)
	pledge.Register(codec)
	transfer.Register(codec)
	project.Register(codec)
	coomunity_project.Register(codec)
	milestone.Register(codec)
	task.Register(codec)
	pledge_task.Register(codec)
	completed_task.Register(codec)
	feature.Register(codec)
	active_request.Register(codec)
	completed.Register(codec)
	active_vote.Register(codec)
	validator.Register(codec)
	link.Register(codec)
	node.Register(codec)
	category.Register(codec)

	// replace:
	entity.Replace(codec)
	request.Replace(codec)
	active_request.Replace(codec)
	completed.Replace(codec)
	vote.Replace(codec)
	active_vote.Replace(codec)

	cdc = codec
}
