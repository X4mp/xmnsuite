package applications

import (
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

// RetrieveValidators is a func that retrieve validators
type RetrieveValidators func(ds datastore.DataStore) ([]Validator, error)

// InfoRequest represents an info request
type InfoRequest interface {
	Version() string
}

// InfoResponse represents an info response
type InfoResponse interface {
	Version() string
	State() State
}

// CommitResponse represents a commit response
type CommitResponse interface {
	AppHash() []byte
	PrevAppHash() []byte
	BlockHeight() int64
}

// Validator represents a validator
type Validator interface {
	IP() net.IP
	PubKey() crypto.PubKey
	Power() int64
}

// Application represents an application
type Application interface {
	GetBlockIndex() int64
	FromBlockIndex() int64
	ToBlockIndex() int64
	Validators() ([]Validator, error)
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
	IP() string
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

// State represents a state
type State interface {
	Hash() []byte
	Height() int64
	Size() int64
	Increment() int64
	Version() string
}

// Database represents the database
type Database interface {
	State(version string) State
	Update(version string) (State, error)
	DataStore() datastore.StoredDataStore
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
	Namespace          string
	Name               string
	ID                 *uuid.UUID
	DirPath            string
	FromBlockIndex     int64
	ToBlockIndex       int64
	Version            string
	Store              datastore.StoredDataStore
	RouterParams       routers.CreateRouterParams
	RetrieveValidators RetrieveValidators
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

// CreateValidatorParams represents the CreateValidator params
type CreateValidatorParams struct {
	IP     net.IP
	PubKey crypto.PubKey
	Power  int64
}

// SDKFunc represents the applications SDK func
var SDKFunc = struct {
	CreateValidator                 func(params CreateValidatorParams) Validator
	CreateInfoRequest               func(params CreateInfoRequestParams) InfoRequest
	CreateApplication               func(params CreateApplicationParams) Application
	CreateApplications              func(params CreateApplicationsParams) Applications
	CreateClientTransactionResponse func(params CreateClientTransactionResponseParams) ClientTransactionResponse
}{
	CreateValidator: func(params CreateValidatorParams) Validator {
		out := createValidator(params.IP, params.PubKey, params.Power)
		return out
	},
	CreateInfoRequest: func(params CreateInfoRequestParams) InfoRequest {
		out := createInfoRequest(params.Version)
		return out
	},
	CreateApplication: func(params CreateApplicationParams) Application {
		//create the router:
		rter := routers.SDKFunc.CreateRouter(params.RouterParams)

		// set some constant:
		stateKey := "state-key"

		// create the database:
		db, dbErr := retrieveOrCreateState(params.Version, stateKey, params.Store)
		if dbErr != nil {
			panic(dbErr)
		}

		//create the application:
		app, appErr := createApplication(params.FromBlockIndex, params.ToBlockIndex, params.Version, db, rter, params.RetrieveValidators)
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
