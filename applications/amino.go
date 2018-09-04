package applications

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	// XMNSuiteApplicationsResourcePointer represents the xmnsuite ResourcePointer resource
	XMNSuiteApplicationsResourcePointer = "xmnsuite/ResourcePointer"

	// XMNSuiteApplicationsResource represents the xmnsuite Resource resource
	XMNSuiteApplicationsResource = "xmnsuite/Resource"

	// XMNSuiteApplicationsInfoRequest represents the xmnsuite InfoRequest resource
	XMNSuiteApplicationsInfoRequest = "xmnsuite/InfoRequest"

	// XMNSuiteApplicationsInfoResponse represents the xmnsuite InfoResponse resource
	XMNSuiteApplicationsInfoResponse = "xmnsuite/InfoResponse"

	// XMNSuiteApplicationsTransactionRequest represents the xmnsuite TransactionRequest resource
	XMNSuiteApplicationsTransactionRequest = "xmnsuite/TransactionRequest"

	// XMNSuiteApplicationsTransactionResponse represents the xmnsuite TransactionResponse resource
	XMNSuiteApplicationsTransactionResponse = "xmnsuite/TransactionResponse"

	// XMNSuiteApplicationsCommitResponse represents the xmnsuite CommitResponse resource
	XMNSuiteApplicationsCommitResponse = "xmnsuite/CommitResponse"

	// XMNSuiteApplicationsQueryRequest represents the xmnsuite QueryRequest resource
	XMNSuiteApplicationsQueryRequest = "xmnsuite/QueryRequest"

	// XMNSuiteApplicationsQueryResponse represents the xmnsuite QueryResponse resource
	XMNSuiteApplicationsQueryResponse = "xmnsuite/QueryResponse"

	// XMNSuiteApplicationsClientTransactionResponse represents the xmnsuite ClientTransactionResponse resource
	XMNSuiteApplicationsClientTransactionResponse = "xmnsuite/ClientTransactionResponse"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// PublicKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)
	}()

	// ResourcePointer
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*ResourcePointer)(nil), nil)
		codec.RegisterConcrete(&resourcePointer{}, XMNSuiteApplicationsResourcePointer, nil)
	}()

	// Resource
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Resource)(nil), nil)
		codec.RegisterConcrete(&resource{}, XMNSuiteApplicationsResource, nil)
	}()

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

	// TransactionRequest
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*TransactionRequest)(nil), nil)
		codec.RegisterConcrete(&transactionRequest{}, XMNSuiteApplicationsTransactionRequest, nil)
	}()

	// TransactionResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*TransactionResponse)(nil), nil)
		codec.RegisterConcrete(&transactionResponse{}, XMNSuiteApplicationsTransactionResponse, nil)
	}()

	// CommitResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*CommitResponse)(nil), nil)
		codec.RegisterConcrete(&commitResponse{}, XMNSuiteApplicationsCommitResponse, nil)
	}()

	// QueryRequest
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*QueryRequest)(nil), nil)
		codec.RegisterConcrete(&queryRequest{}, XMNSuiteApplicationsQueryRequest, nil)
	}()

	// QueryResponse
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*QueryResponse)(nil), nil)
		codec.RegisterConcrete(&queryResponse{}, XMNSuiteApplicationsQueryResponse, nil)
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
