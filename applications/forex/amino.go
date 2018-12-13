package forex

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/fiatchain"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	bank.Register(codec)
	category.Register(codec)
	currency.Register(codec)
	deposit.Register(codec)
	fiatchain.Register(codec)
}
