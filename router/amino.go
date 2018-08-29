package router

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var cdc = amino.NewCodec()

// RequestAminoRoute represents the Request amino route
const RequestAminoRoute = "datamint/router/Request"

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
		codec.RegisterInterface((*TrxChkResponse)(nil), nil)
		codec.RegisterConcrete(trxChkResponse{}, TrxChkResponseAminoRoute, nil)
	}()

	// PublicKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)
	}()

	// PrivateKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PrivKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.Ed25519PrivKeyAminoRoute, nil)
	}()
}
