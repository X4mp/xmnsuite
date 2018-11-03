package request

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/framework/user"
)

const (

	// XMNSuiteBlockchainsFrameworkVoteRequest represents the xmnsuite framework vote Request
	XMNSuiteBlockchainsFrameworkVoteRequest = "xmnsuite/blockchains/framework/vote/Request"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Request
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Request)(nil), nil)
		codec.RegisterConcrete(&request{}, XMNSuiteBlockchainsFrameworkVoteRequest, nil)
	}()
}
