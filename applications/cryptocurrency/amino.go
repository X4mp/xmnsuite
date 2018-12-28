package cryptocurrency

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/chain"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/web"
	core "github.com/xmnservices/xmnsuite/blockchains/core"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	address.Register(codec)
	chain.Register(codec)
	deposit.Register(codec)
	offer.Register(codec)
	web.Register(codec)
	core.Register(codec)

	cdc = codec
}
