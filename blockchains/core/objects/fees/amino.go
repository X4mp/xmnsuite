package fees

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

const (

	// XMNSuiteBlockchainsCoreFee represents the xmnsuite core Fee
	XMNSuiteBlockchainsCoreFee = "xmnsuite/blockchains/core/Fee"

	// XMNSuiteBlockchainsCoreNormalizedFee represents the xmnsuite core NormalizedFee
	XMNSuiteBlockchainsCoreNormalizedFee = "xmnsuite/blockchains/core/NormalizedFee"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	user.Register(codec)
	keyname.Register(codec)

	// Fee
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Fee)(nil), nil)
		codec.RegisterConcrete(&fee{}, XMNSuiteBlockchainsCoreFee, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&normalizedFee{}, XMNSuiteBlockchainsCoreNormalizedFee, nil)
	}()
}

// Replace replaces the amino codec
func Replace(codec *amino.Codec) {
	// replace:
	cdc = codec

	// register again:
	Register(cdc)
}
