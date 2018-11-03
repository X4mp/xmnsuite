package genesis

import (
	amino "github.com/tendermint/go-amino"
	deposit "github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
)

const (

	// XMNSuiteApplicationsXMNGenesis represents the xmnsuite xmn Genesis resource
	XMNSuiteApplicationsXMNGenesis = "xmnsuite/xmn/Genesis"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	deposit.Register(codec)
	token.Register(codec)

	// Genesis
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Genesis)(nil), nil)
		codec.RegisterConcrete(&genesis{}, XMNSuiteApplicationsXMNGenesis, nil)
	}()
}
