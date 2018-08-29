package applications

import (
	datastore "github.com/XMNBlockchain/datamint/datastore"
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

// CreateInfoRequestParams represents the CreateInfoRequest params
type CreateInfoRequestParams struct {
	Version string
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateInfoRequest func(params CreateInfoRequestParams) InfoRequest
}{
	CreateInfoRequest: func(params CreateInfoRequestParams) InfoRequest {
		out := createInfoRequest(params.Version)
		return out
	},
}
