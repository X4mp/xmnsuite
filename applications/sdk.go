package applications

import (
	"errors"

	datastore "github.com/XMNBlockchain/datamint/datastore"
	objects "github.com/XMNBlockchain/datamint/objects"
	crypto "github.com/tendermint/tendermint/crypto"
)

// SaveTransactionFn represents a save transaction func
type SaveTransactionFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) TransactionResponse

// DeleteTransactionFn represents a delete transaction func
type DeleteTransactionFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) TransactionResponse

// QueryFn represents a query func.  The return values are: code, key, value, log
type QueryFn func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) QueryResponse

const (
	// IsSuccessful represents a successful query and/or transaction
	IsSuccessful = iota

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
	Info(req InfoRequest) InfoResponse
	Transact(req TransactionRequest) TransactionResponse
	CheckTransact(req TransactionRequest) TransactionResponse
	Commit() CommitResponse
	Query(req QueryRequest) QueryResponse
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

// CreateQueryRequestParams represents the CreateQueryRequest params
type CreateQueryRequestParams struct {
	Ptr    ResourcePointer
	Sig    []byte
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
	Version      string
	DataStore    datastore.DataStore
	RouterParams CreateRouterParams
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateResourcePointer    func(params CreateResourcePointerParams) ResourcePointer
	CreateResource           func(params CreateResourceParams) Resource
	CreateInfoRequest        func(params CreateInfoRequestParams) InfoRequest
	CreateTransactionRequest func(params CreateTransactionRequestParams) TransactionRequest
	CreateQueryRequest       func(params CreateQueryRequestParams) QueryRequest
	CreateApplication        func(params CreateApplicationParams) Application
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
	CreateApplication: func(params CreateApplicationParams) Application {

		// fetch the roles and users:
		rols := params.DataStore.Roles()
		usrs := params.DataStore.Users()

		//create the routes:
		rtes := map[int][]Route{}
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
		app, appErr := createApplication(params.Version, stateKey, storedState, params.DataStore, rter)
		if appErr != nil {
			panic(appErr)
		}

		return app
	},
}
