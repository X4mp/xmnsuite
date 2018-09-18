package routers

import (
	"errors"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
)

// SaveTransactionFn represents a save transaction func
type SaveTransactionFn func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (TransactionResponse, error)

// DeleteTransactionFn represents a delete transaction func
type DeleteTransactionFn func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (TransactionResponse, error)

// QueryFn represents a query func.  The return values are: code, key, value, log
type QueryFn func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error)

const (
	// IsSuccessful represents a successful query and/or transaction
	IsSuccessful = iota

	// NotFound represents a resource not found
	NotFound

	// ServerError represents a server error
	ServerError

	// IsUnAuthorized represents an un-authorized query and/or transaction
	IsUnAuthorized

	// IsUnAuthenticated represents an un-authenticated query and/or transaction
	IsUnAuthenticated

	// RouteNotFound represents a route not found query and/or transaction
	RouteNotFound

	// InvalidRoute represents an invalid route
	InvalidRoute

	// InvalidRequest represents an invalid request
	InvalidRequest
)

const (
	// Save represents the save transaction method
	Save = iota

	// Delete represents the delete transaction method
	Delete

	// Retrieve represents the retrieve query method
	Retrieve
)

// ResourcePointer represents a resource pointer
type ResourcePointer interface {
	From() crypto.PublicKey
	Path() string
	Hash() string
}

// Resource represents a resource
type Resource interface {
	Pointer() ResourcePointer
	Data() []byte
	Hash() string
}

// TransactionRequest represents a transaction request
type TransactionRequest interface {
	Resource() Resource
	Pointer() ResourcePointer
	Signature() crypto.Signature
}

// TransactionResponse represents a transaction response
type TransactionResponse interface {
	Code() int
	Log() string
	GazUsed() int64
	Tags() map[string][]byte
}

// QueryRequest represents a query request
type QueryRequest interface {
	Pointer() ResourcePointer
	Signature() crypto.Signature
}

// QueryResponse represents a query response
type QueryResponse interface {
	Code() int
	Log() string
	Key() string
	Value() []byte
}

// Handler represents a router handler
type Handler interface {
	SaveTransaction() SaveTransactionFn
	DeleteTransaction() DeleteTransactionFn
	Query() QueryFn
	IsWrite() bool
}

// PreparedHandler represents a prepated handler
type PreparedHandler interface {
	Path() string
	Params() map[string]string
	Handler() Handler
}

// Route represents a route
type Route interface {
	Matches(from crypto.PublicKey, path string) bool
	Handler(from crypto.PublicKey, path string) PreparedHandler
}

// Router represents a router
type Router interface {
	Route(from crypto.PublicKey, path string, method int) PreparedHandler
}

/*
 * SDK Params
 */

// CreateResourcePointerParams represents the CreateResourcePointer params
type CreateResourcePointerParams struct {
	From crypto.PublicKey
	Path string
}

// CreateResourceParams represents the CreateResource params
type CreateResourceParams struct {
	ResPtr ResourcePointer
	Data   []byte
}

// CreateTransactionRequestParams represents the CreateTransactionRequest params
type CreateTransactionRequestParams struct {
	Res    Resource
	Ptr    ResourcePointer
	Sig    crypto.Signature
	JSData []byte
}

// CreateTransactionResponseParams represents the CreateTransactionResponse params
type CreateTransactionResponseParams struct {
	Code    int
	Log     string
	GazUsed int64
	Tags    map[string][]byte
}

// CreateQueryRequestParams represents the CreateQueryRequest params
type CreateQueryRequestParams struct {
	Ptr    ResourcePointer
	Sig    crypto.Signature
	JSData []byte
}

// CreateQueryResponseParams represents the CreateQueryResponse params
type CreateQueryResponseParams struct {
	Code   int
	Log    string
	Key    string
	Value  []byte
	JSData []byte
}

// CreateRouteParams represents the CreateRoute params
type CreateRouteParams struct {
	Pattern  string
	SaveTrx  SaveTransactionFn
	DelTrx   DeleteTransactionFn
	QueryTrx QueryFn
}

// CreateRouterParams represents the CreateRouter params
type CreateRouterParams struct {
	DataStore  datastore.DataStore
	RoleKey    string
	RtesParams []CreateRouteParams
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateResourcePointer     func(params CreateResourcePointerParams) ResourcePointer
	CreateResource            func(params CreateResourceParams) Resource
	CreateTransactionRequest  func(params CreateTransactionRequestParams) TransactionRequest
	CreateTransactionResponse func(params CreateTransactionResponseParams) TransactionResponse
	CreateQueryRequest        func(params CreateQueryRequestParams) QueryRequest
	CreateQueryResponse       func(params CreateQueryResponseParams) QueryResponse
	CreateRouter              func(params CreateRouterParams) Router
}{
	CreateResourcePointer: func(params CreateResourcePointerParams) ResourcePointer {
		out := createResourcePointer(params.From, params.Path)
		return out
	},
	CreateResource: func(params CreateResourceParams) Resource {
		out := createResource(params.ResPtr, params.Data)
		return out
	},
	CreateTransactionRequest: func(params CreateTransactionRequestParams) TransactionRequest {
		if params.JSData != nil {
			out := new(transactionRequest)
			jsErr := cdc.UnmarshalJSON(params.JSData, out)
			if jsErr != nil {
				panic(jsErr)
			}

			return out
		}

		if params.Ptr != nil {
			out, outErr := createTransactionRequestWithResourcePointer(params.Ptr, params.Sig)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		if params.Res != nil {
			out, outErr := createTransactionRequestWithResource(params.Res, params.Sig)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		panic(errors.New("the params must contain a Resource or a PointerResource"))
	},
	CreateTransactionResponse: func(params CreateTransactionResponseParams) TransactionResponse {
		if params.GazUsed != 0 && params.Tags != nil {
			out, outErr := createTransactionResponse(params.Code, params.Log, params.GazUsed, params.Tags)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := createFreeTransactionResponse(params.Code, params.Log)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateQueryRequest: func(params CreateQueryRequestParams) QueryRequest {
		if params.JSData != nil {
			qr := new(queryRequest)
			jsErr := cdc.UnmarshalJSON(params.JSData, qr)
			if jsErr != nil {
				panic(jsErr)
			}

			return qr
		}

		out, outErr := createQueryRequest(params.Ptr, params.Sig)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateQueryResponse: func(params CreateQueryResponseParams) QueryResponse {
		if params.JSData != nil {
			out := new(queryResponse)
			jsErr := cdc.UnmarshalJSON(params.JSData, out)
			if jsErr != nil {
				panic(jsErr)
			}

			return out
		}

		if params.Key == "" && params.Value == nil {
			out, outErr := createEmptyQueryResponse(params.Code, params.Log)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := createQueryResponse(params.Code, params.Log, params.Key, params.Value)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateRouter: func(params CreateRouterParams) Router {
		rtes := map[int][]Route{}
		rols := params.DataStore.Roles()
		usrs := params.DataStore.Users()
		for _, oneRteParams := range params.RtesParams {
			//create handler:
			handlr, rteType, handlrErr := createHandler(oneRteParams.SaveTrx, oneRteParams.DelTrx, oneRteParams.QueryTrx)
			if handlrErr != nil {
				panic(handlrErr)
			}

			//create route:
			rte, rteErr := createRoute(params.RoleKey, rols, usrs, oneRteParams.Pattern, handlr)
			if rteErr != nil {
				panic(rteErr)
			}

			// init the list for the given route type:
			if _, ok := rtes[rteType]; !ok {
				rtes[rteType] = []Route{}
			}

			// add the route:
			rtes[rteType] = append(rtes[rteType], rte)
		}

		//create the router:
		rter := createRouter(rtes)
		return rter
	},
}
