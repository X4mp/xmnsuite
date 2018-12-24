package forex

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/fiatchain"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
	entity "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	core.Register(codec)
	bank.Register(codec)
	category.Register(codec)
	currency.Register(codec)
	deposit.Register(codec)
	fiatchain.Register(codec)

	// replace:
	entity.Replace(cdc)
	request.Replace(cdc)
	vote.Replace(cdc)
}
