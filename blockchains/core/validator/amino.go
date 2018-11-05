package validator

import (
	amino "github.com/tendermint/go-amino"
)

const (

	// XMNSuiteApplicationsXMNValidator represents the xmnsuite xmn Validator resource
	XMNSuiteApplicationsXMNValidator = "xmnsuite/xmn/Validator"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Validator
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Validator)(nil), nil)
		codec.RegisterConcrete(&validator{}, XMNSuiteApplicationsXMNValidator, nil)
	}()
}
