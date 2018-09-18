package applications

import (
	datastore "github.com/xmnservices/xmnsuite/datastore"
	objects "github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/routers"
)

// InfoRequest represents an info request
type InfoRequest interface {
	Version() string
}

// InfoResponse represents an info response
type InfoResponse interface {
	Size() int64
	Version() string
}

// CommitResponse represents a commit response
type CommitResponse interface {
	AppHash() []byte
	BlockHeight() int64
}

// Application represents an application
type Application interface {
	GetBlockIndex() int64
	FromBlockIndex() int64
	ToBlockIndex() int64
	Info(req InfoRequest) InfoResponse
	Transact(req routers.TransactionRequest) routers.TransactionResponse
	CheckTransact(req routers.TransactionRequest) routers.TransactionResponse
	Commit() CommitResponse
	Query(req routers.QueryRequest) routers.QueryResponse
}

// Applications represents an application
type Applications interface {
	RetrieveBlockIndex() int64
	RetrieveByBlockIndex(blkIndex int64) (Application, error)
}

// ClientTransactionResponse represents a client transaction response
type ClientTransactionResponse interface {
	Check() routers.TransactionResponse
	Transaction() routers.TransactionResponse
	Height() int64
	Hash() []byte
}

// Client represents an application client
type Client interface {
	Query(req routers.QueryRequest) (routers.QueryResponse, error)
	Transact(req routers.TransactionRequest) (ClientTransactionResponse, error)
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

// CreateInfoRequestParams represents the CreateInfoRequest params
type CreateInfoRequestParams struct {
	Version string
}

// CreateApplicationParams represents the CreateApplication params
type CreateApplicationParams struct {
	FromBlockIndex int64
	ToBlockIndex   int64
	Version        string
	DataStore      datastore.DataStore
	RouterParams   routers.CreateRouterParams
}

// CreateApplicationsParams represents the CreateApplications params
type CreateApplicationsParams struct {
	Apps []Application
}

// CreateClientTransactionResponseParams represents the CreateClientTransactionResponse params
type CreateClientTransactionResponseParams struct {
	Chk    routers.TransactionResponse
	Trx    routers.TransactionResponse
	Height int64
	Hash   []byte
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateInfoRequest               func(params CreateInfoRequestParams) InfoRequest
	CreateApplication               func(params CreateApplicationParams) Application
	CreateApplications              func(params CreateApplicationsParams) Applications
	CreateClientTransactionResponse func(params CreateClientTransactionResponseParams) ClientTransactionResponse
}{
	CreateInfoRequest: func(params CreateInfoRequestParams) InfoRequest {
		out := createInfoRequest(params.Version)
		return out
	},
	CreateApplication: func(params CreateApplicationParams) Application {
		//create the router:
		rter := routers.SDKFunc.CreateRouter(params.RouterParams)

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
