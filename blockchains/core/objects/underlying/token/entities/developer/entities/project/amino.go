package project

import (
	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

const (

	// XMNSuiteApplicationsXMNProject represents the xmnsuite xmn Project resource
	XMNSuiteApplicationsXMNProject = "xmnsuite/xmn/Project"

	// XMNSuiteApplicationsXMNNormalizedProject represents the xmnsuite xmn Normalized Project resource
	XMNSuiteApplicationsXMNNormalizedProject = "xmnsuite/xmn/Normalized/Project"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// dependencies:
	user.Register(codec)

	// Project
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Project)(nil), nil)
		codec.RegisterConcrete(&project{}, XMNSuiteApplicationsXMNProject, nil)
	}()

	// Normalized
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Normalized)(nil), nil)
		codec.RegisterConcrete(&storableProject{}, XMNSuiteApplicationsXMNNormalizedProject, nil)
	}()
}
