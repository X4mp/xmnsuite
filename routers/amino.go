package routers

import (
	amino "github.com/tendermint/go-amino"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

const (
	// XMNSuiteApplicationsResourcePointer represents the xmnsuite ResourcePointer resource
	XMNSuiteApplicationsResourcePointer = "xmnsuite/ResourcePointer"

	// XMNSuiteApplicationsResource represents the xmnsuite Resource resource
	XMNSuiteApplicationsResource = "xmnsuite/Resource"

	// XMNSuiteApplicationsTransactionRequest represents the xmnsuite TransactionRequest resource
	XMNSuiteApplicationsTransactionRequest = "xmnsuite/TransactionRequest"

	// XMNSuiteApplicationsTransactionResponse represents the xmnsuite TransactionResponse resource
	XMNSuiteApplicationsTransactionResponse = "xmnsuite/TransactionResponse"

	// XMNSuiteApplicationsQueryRequest represents the xmnsuite QueryRequest resource
	XMNSuiteApplicationsQueryRequest = "xmnsuite/QueryRequest"

	// XMNSuiteApplicationsQueryResponse represents the xmnsuite QueryResponse resource
	XMNSuiteApplicationsQueryResponse = "xmnsuite/QueryResponse"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// Dependencies
	crypto.Register(codec)

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
}
