package applications

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	routers "github.com/xmnservices/xmnsuite/routers"
)

const (

	// XMNSuiteApplicationsInfoRequest represents the xmnsuite InfoRequest resource
	XMNSuiteApplicationsInfoRequest = "xmnsuite/InfoRequest"

	// XMNSuiteApplicationsInfoResponse represents the xmnsuite InfoResponse resource
	XMNSuiteApplicationsInfoResponse = "xmnsuite/InfoResponse"

	// XMNSuiteApplicationsCommitResponse represents the xmnsuite CommitResponse resource
	XMNSuiteApplicationsCommitResponse = "xmnsuite/CommitResponse"

	// XMNSuiteApplicationsClientTransactionResponse represents the xmnsuite ClientTransactionResponse resource
	XMNSuiteApplicationsClientTransactionResponse = "xmnsuite/ClientTransactionResponse"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	crypto.Register(codec)
	routers.Register(codec)

	// InfoRequest
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*InfoRequest)(nil), nil)
		codec.RegisterConcrete(&infoRequest{}, XMNSuiteApplicationsInfoRequest, nil)
	}()

	// InfoResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*InfoResponse)(nil), nil)
		codec.RegisterConcrete(&infoResponse{}, XMNSuiteApplicationsInfoResponse, nil)
	}()

	// CommitResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*CommitResponse)(nil), nil)
		codec.RegisterConcrete(&commitResponse{}, XMNSuiteApplicationsCommitResponse, nil)
	}()

	// ClientTransactionResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*ClientTransactionResponse)(nil), nil)
		codec.RegisterConcrete(&clientTransactionResponse{}, XMNSuiteApplicationsClientTransactionResponse, nil)
	}()
}
