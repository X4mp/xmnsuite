package router

import (
	"errors"
	"fmt"

	crypto "github.com/tendermint/tendermint/crypto"
)

// QueryHandlerFn represents a query handler func
type QueryHandlerFn func(req Request) QueryResponse

// TrxHandlerFn represents a transaction handler func
type TrxHandlerFn func(req Request) TrxResponse

// TrxChkHandlerFn represents a transaction check handler func
type TrxChkHandlerFn func(req TrxChkRequest) TrxChkResponse

// Request represents a request
type Request interface {
	From() crypto.PubKey
	Path() string
	Data() []byte
	Signature() []byte
}

// QueryResponse represents a query response
type QueryResponse interface {
	IsSuccess() bool
	IsAuthorized() bool
	HasInsufficientFunds() bool
	GazUsed() int64
	Log() string
	Data() []byte
	UnMarshal(v interface{}) error
}

// TrxResponse represents a transaction response
type TrxResponse interface {
	IsSuccess() bool
	IsAuthorized() bool
	HasInsufficientFunds() bool
	Tags() map[string][]byte
	GazUsed() int64
	Log() string
}

// TrxChkRequest represents a transaction check request
type TrxChkRequest interface {
	From() crypto.PubKey
	Path() string
	DataSizeInBytes() int64
	Signature() []byte
}

// TrxChkResponse represents a transaction check response
type TrxChkResponse interface {
	CanBeExecuted() bool
	CanBeAuthorized() bool
	GazWanted() int64
	Log() string
}

// QueryRoute represents a query route
type QueryRoute interface {
	Matches(req Request) bool
	Handler() QueryHandlerFn
}

// TrxChkRoute represents a transaction check route
type TrxChkRoute interface {
	Matches(req TrxChkRequest) bool
	Handler() TrxChkHandlerFn
}

// TrxRoute represents a transaction route
type TrxRoute interface {
	Matches(req Request) bool
	Handler() TrxHandlerFn
}

// Router represents the router
type Router interface {
	Start() error
	Stop()
	Query(request Request) QueryResponse
	CheckTrx(request TrxChkRequest) TrxChkResponse
	Transact(request Request) TrxResponse
}

/*
 * Params
 */

// CreateRequestParams represents the CreateRequest params
type CreateRequestParams struct {
	From      crypto.PubKey
	Path      string
	Data      []byte
	Signature []byte
	JSData    []byte
}

// CreateTrxChkRequestParams represents the CreateTrxChkRequest params
type CreateTrxChkRequestParams struct {
	From            crypto.PubKey
	Path            string
	DataSizeInBytes int64
	Signature       []byte
	JSData          []byte
}

// CreateRouterParams represents the CreateRouter params
type CreateRouterParams struct {
	QueryRtes  []QueryRoute
	TrxChkRtes []TrxChkRoute
	TrxRtes    []TrxRoute
}

// CreateQueryResponseParams represents the CreateQueryResponse params
type CreateQueryResponseParams struct {
	IsSuccess    bool
	IsAuthorized bool
	IsNFS        bool
	GazUsed      int64
	Data         []byte
	Log          string
	JSData       []byte
}

// CreateTrxChkResponseParams represents the CreateTrxChkResponse params
type CreateTrxChkResponseParams struct {
	CanBeExecuted   bool
	CanBeAuthorized bool
	GazWanted       int64
	Log             string
	JSData          []byte
}

// CreateTrxResponseParams represents the CreateTrxResponse params
type CreateTrxResponseParams struct {
	IsSuccess    bool
	IsAuthorized bool
	IsNFS        bool
	Tags         map[string][]byte
	GazUsed      int64
	Log          string
	JSData       []byte
}

// SDKFunc represents the router SDK func
var SDKFunc = struct {
	CreateRequest        func(params CreateRequestParams) Request
	CreateTrxChkRequest  func(params CreateTrxChkRequestParams) TrxChkRequest
	CreateQueryResponse  func(param CreateQueryResponseParams) QueryResponse
	CreateTrxChkResponse func(param CreateTrxChkResponseParams) TrxChkResponse
	CreateTrxResponse    func(param CreateTrxResponseParams) TrxResponse
	CreateRouter         func(params CreateRouterParams) Router
}{
	CreateRequest: func(params CreateRequestParams) Request {
		if params.JSData != nil {
			ptr := new(request)
			err := cdc.UnmarshalJSON(params.JSData, ptr)
			if err != nil {
				str := fmt.Sprintf("invalid json data: %s", err.Error())
				panic(errors.New(str))
			}

			return ptr
		}

		out, outErr := createRequest(params.From, params.Path, params.Data, params.Signature)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateTrxChkRequest: func(params CreateTrxChkRequestParams) TrxChkRequest {
		if params.JSData != nil {
			ptr := new(trxChkRequest)
			err := cdc.UnmarshalJSON(params.JSData, ptr)
			if err != nil {
				str := fmt.Sprintf("invalid json data: %s", err.Error())
				panic(errors.New(str))
			}

			return ptr
		}

		out, outErr := createTrxChkRequest(params.From, params.Path, params.DataSizeInBytes, params.Signature)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateQueryResponse: func(params CreateQueryResponseParams) QueryResponse {
		if params.JSData != nil {
			ptr := new(queryResponse)
			err := cdc.UnmarshalJSON(params.JSData, ptr)
			if err != nil {
				str := fmt.Sprintf("invalid json data: %s", err.Error())
				panic(errors.New(str))
			}

			return ptr
		}

		out, outErr := createQueryResponse(params.IsSuccess, params.IsAuthorized, params.IsNFS, params.GazUsed, params.Log, params.Data)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateTrxChkResponse: func(params CreateTrxChkResponseParams) TrxChkResponse {
		if params.JSData != nil {
			ptr := new(trxChkResponse)
			err := cdc.UnmarshalJSON(params.JSData, ptr)
			if err != nil {
				str := fmt.Sprintf("invalid json data: %s", err.Error())
				panic(errors.New(str))
			}

			return ptr
		}

		out, outErr := createTrxChkResponse(params.CanBeExecuted, params.CanBeAuthorized, params.GazWanted, params.Log)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateTrxResponse: func(params CreateTrxResponseParams) TrxResponse {
		if params.JSData != nil {
			ptr := new(trxResponse)
			err := cdc.UnmarshalJSON(params.JSData, ptr)
			if err != nil {
				str := fmt.Sprintf("invalid json data: %s", err.Error())
				panic(errors.New(str))
			}

			return ptr
		}

		out, outErr := createTrxResponse(params.IsSuccess, params.IsAuthorized, params.IsNFS, params.Tags, params.GazUsed, params.Log)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateRouter: func(params CreateRouterParams) Router {
		return createRouter(params.QueryRtes, params.TrxChkRtes, params.TrxRtes)
	},
}
