package xmn

import (
	amino "github.com/tendermint/go-amino"
	genesis "github.com/xmnservices/xmnsuite/blockchains/xmn/genesis"
	initial_deposit "github.com/xmnservices/xmnsuite/blockchains/xmn/initial_deposit"
	token "github.com/xmnservices/xmnsuite/blockchains/xmn/token"
	user "github.com/xmnservices/xmnsuite/blockchains/xmn/user"
	wallet "github.com/xmnservices/xmnsuite/blockchains/xmn/wallet"
	applications "github.com/xmnservices/xmnsuite/routers"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	applications.Register(codec)
	wallet.Register(codec)
	user.Register(codec)
	token.Register(codec)
	genesis.Register(codec)
	initial_deposit.Register(codec)
}
