package token

import (
	amino "github.com/tendermint/go-amino"
)

const (

	// XMNSuiteApplicationsXMNToken represents the xmnsuite xmn Token resource
	XMNSuiteApplicationsXMNToken = "xmnsuite/xmn/Token"

	// XMNSuiteApplicationsXMNNormalizedToken represents the xmnsuite xmn Normalized Token resource
	XMNSuiteApplicationsXMNNormalizedToken = "xmnsuite/xmn/Normalized/Token"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Token
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Token)(nil), nil)
		codec.RegisterConcrete(&token{}, XMNSuiteApplicationsXMNToken, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableToken{}, XMNSuiteApplicationsXMNNormalizedToken, nil)
	}()
}
