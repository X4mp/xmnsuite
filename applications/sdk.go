package applications

import (
	"errors"

	crypto "github.com/tendermint/tendermint/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	objects "github.com/xmnservices/xmnsuite/objects"
)

// SaveTransactionFn represents a save transaction func
type SaveTransactionFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) (TransactionResponse, error)

// DeleteTransactionFn represents a delete transaction func
type DeleteTransactionFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) (TransactionResponse, error)

// QueryFn represents a query func.  The return values are: code, key, value, log
type QueryFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) (QueryResponse, error)

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
	From() crypto.PubKey
	Path() string
	Hash() []byte
}

// Resource represents a resource
type Resource interface {
	Pointer() ResourcePointer
	Data() []byte
	Hash() []byte
}

// InfoRequest represents an info request
type InfoRequest interface {
	Version() string
}

// InfoResponse represents an info response
type InfoResponse interface {
	Size() int64
	Version() string
	LastBlockHeight() int64
	LastBlockAppHash() []byte
}

// TransactionRequest represents a transaction request
type TransactionRequest interface {
	Resource() Resource
	Pointer() ResourcePointer
	Signature() []byte
}

// TransactionResponse represents a transaction response
type TransactionResponse interface {
	Code() int
	Log() string
	GazUsed() int64
	Tags() map[string][]byte
}

// CommitResponse represents a commit response
type CommitResponse interface {
	AppHash() []byte
	BlockHeight() int64
}

// QueryRequest represents a query request
type QueryRequest interface {
	Pointer() ResourcePointer
	Signature() []byte
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
	Matches(from crypto.PubKey, path string) bool
	Handler(from crypto.PubKey, path string) PreparedHandler
}

// Router represents a router
type Router interface {
	Route(from crypto.PubKey, path string, method int) PreparedHandler
}

// Application represents an application
type Application interface {
	GetBlockIndex() int64
	FromBlockIndex() int64
	ToBlockIndex() int64
	Info(req InfoRequest) InfoResponse
	Transact(req TransactionRequest) TransactionResponse
	CheckTransact(req TransactionRequest) TransactionResponse
	Commit() CommitResponse
	Query(req QueryRequest) QueryResponse
}

// Applications represents an application
type Applications interface {
	RetrieveBlockIndex() int64
	RetrieveByBlockIndex(blkIndex int64) (Application, error)
}

// ClientTransactionResponse represents a client transaction response
type ClientTransactionResponse interface {
	Check() TransactionResponse
	Transaction() TransactionResponse
	Height() int64
	Hash() []byte
}

// Client represents an application client
type Client interface {
	Query(req QueryRequest) (QueryResponse, error)
	Transact(req TransactionRequest) (ClientTransactionResponse, error)
}

// Node represents a node in which an application is running
type Node interface {
	GetAddress() string
	GetClient() (Client, error)
	Start() error
	Stop() error
}

/*
 * SDK Params
 */

// CreateResourcePointerParams represents the CreateResourcePointer params
type CreateResourcePointerParams struct {
	From crypto.PubKey
	Path string
}

// CreateResourceParams represents the CreateResource params
type CreateResourceParams struct {
	ResPtr ResourcePointer
	Data   []byte
}

// CreateInfoRequestParams represents the CreateInfoRequest params
type CreateInfoRequestParams struct {
	Version string
}

// CreateTransactionRequestParams represents the CreateTransactionRequest params
type CreateTransactionRequestParams struct {
	Res    Resource
	Ptr    ResourcePointer
	Sig    []byte
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
	Sig    []byte
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

// CreateApplicationParams represents the CreateApplication params
type CreateApplicationParams struct {
	FromBlockIndex int64
	ToBlockIndex   int64
	Version        string
	DataStore      datastore.DataStore
	RouterParams   CreateRouterParams
}

// CreateApplicationsParams represents the CreateApplications params
type CreateApplicationsParams struct {
	Apps []Application
}

// CreateClientTransactionResponseParams represents the CreateClientTransactionResponse params
type CreateClientTransactionResponseParams struct {
	Chk    TransactionResponse
	Trx    TransactionResponse
	Height int64
	Hash   []byte
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateResourcePointer           func(params CreateResourcePointerParams) ResourcePointer
	CreateResource                  func(params CreateResourceParams) Resource
	CreateInfoRequest               func(params CreateInfoRequestParams) InfoRequest
	CreateTransactionRequest        func(params CreateTransactionRequestParams) TransactionRequest
	CreateTransactionResponse       func(params CreateTransactionResponseParams) TransactionResponse
	CreateQueryRequest              func(params CreateQueryRequestParams) QueryRequest
	CreateQueryResponse             func(params CreateQueryResponseParams) QueryResponse
	CreateApplication               func(params CreateApplicationParams) Application
	CreateApplications              func(params CreateApplicationsParams) Applications
	CreateClientTransactionResponse func(params CreateClientTransactionResponseParams) ClientTransactionResponse
}{
	CreateResourcePointer: func(params CreateResourcePointerParams) ResourcePointer {
		out := createResourcePointer(params.From, params.Path)
		return out
	},
	CreateResource: func(params CreateResourceParams) Resource {
		out := createResource(params.ResPtr, params.Data)
		return out
	},
	CreateInfoRequest: func(params CreateInfoRequestParams) InfoRequest {
		out := createInfoRequest(params.Version)
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
	CreateApplication: func(params CreateApplicationParams) Application {
		//create the routes:
		rtes := map[int][]Route{}
		routerDS := params.RouterParams.DataStore
		rols := routerDS.Roles()
		usrs := routerDS.Users()
		for _, oneRteParams := range params.RouterParams.RtesParams {
			//create handler:
			handlr, rteType, handlrErr := createHandler(oneRteParams.SaveTrx, oneRteParams.DelTrx, oneRteParams.QueryTrx)
			if handlrErr != nil {
				panic(handlrErr)
			}

			//create route:
			rte, rteErr := createRoute(params.RouterParams.RoleKey, rols, usrs, oneRteParams.Pattern, handlr)
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

		// set some constant:
		stateKey := "state-key"

		// create/retrieve the stored state:
		stateObjects := objects.SDKFunc.Create()
		storedState, storedStateErr := retrieveOrCreateState(params.Version, stateKey, stateObjects)
		if storedStateErr != nil {
			panic(storedStateErr)
		}

		//create the application:
		app, appErr := createApplication(params.FromBlockIndex, params.ToBlockIndex, params.Version, stateKey, storedState, params.DataStore, rter)
		if appErr != nil {
			panic(appErr)
		}

		return app
	},
	CreateApplications: func(params CreateApplicationsParams) Applications {
		out := createApplications(params.Apps)
		return out
	},
	CreateClientTransactionResponse: func(params CreateClientTransactionResponseParams) ClientTransactionResponse {
		out := createClientTransactionResponse(params.Chk, params.Trx, params.Height, params.Hash)
		return out
	},
}
