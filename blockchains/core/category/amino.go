package category

import (
	amino "github.com/tendermint/go-amino"
)

const (

	// XMNSuiteApplicationsXMNCategory represents the xmnsuite xmn Category resource
	XMNSuiteApplicationsXMNCategory = "xmnsuite/xmn/Category"

	// XMNSuiteApplicationsXMNNormalizedCategory represents the xmnsuite xmn NormalizedCategory resource
	XMNSuiteApplicationsXMNNormalizedCategory = "xmnsuite/xmn/NormalizedCategory"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Category
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Category)(nil), nil)
		codec.RegisterConcrete(&category{}, XMNSuiteApplicationsXMNCategory, nil)
	}()

	// normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableCategory{}, XMNSuiteApplicationsXMNNormalizedCategory, nil)
	}()
}
