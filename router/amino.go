package router

import (
	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

// RequestAminoRoute represents the Request amino route
const RequestAminoRoute = "datamint/router/Request"

// TrxChkRequestAminoRoute represents the TrxChkRequest amino route
const TrxChkRequestAminoRoute = "datamint/router/TrxChkRequest"

// QueryResponseAminoRoute represents the QueryResponse amino route
const QueryResponseAminoRoute = "datamint/router/QueryResponse"

// TrxResponseAminoRoute represents the TrxResponse amino route
const TrxResponseAminoRoute = "datamint/router/TrxResponse"

// TrxChkResponseAminoRoute represents the TrxChkResponse amino route
const TrxChkResponseAminoRoute = "datamint/router/TrxChkResponse"

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*Request)(nil), nil)
		codec.RegisterConcrete(request{}, RequestAminoRoute, nil)
	}()

	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*QueryResponse)(nil), nil)
		codec.RegisterConcrete(queryResponse{}, QueryResponseAminoRoute, nil)
	}()

	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*TrxResponse)(nil), nil)
		codec.RegisterConcrete(trxResponse{}, TrxResponseAminoRoute, nil)
	}()

	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*TrxChkRequest)(nil), nil)
		codec.RegisterConcrete(trxChkRequest{}, TrxChkRequestAminoRoute, nil)
	}()

	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*TrxChkResponse)(nil), nil)
		codec.RegisterConcrete(trxChkResponse{}, TrxChkResponseAminoRoute, nil)
	}()
}
