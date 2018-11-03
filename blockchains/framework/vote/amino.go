package vote

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/framework/user"
)

const (

	// XMNSuiteBlockchainsFrameworkVote represents the xmnsuite framework vote Vote
	XMNSuiteBlockchainsFrameworkVote = "xmnsuite/blockchains/framework/vote/Vote"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)

	// Vote
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Vote)(nil), nil)
		codec.RegisterConcrete(&vote{}, XMNSuiteBlockchainsFrameworkVote, nil)
	}()
}
