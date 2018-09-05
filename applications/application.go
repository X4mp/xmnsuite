package applications

import (
	"encoding/json"
	"errors"
	"fmt"

	datastore "github.com/XMNBlockchain/datamint/datastore"
)

/*
 * Application
 */

type application struct {
	fromIndex   int64
	toIndex     int64
	version     string
	stateKey    string
	router      Router
	store       datastore.DataStore
	storedState *storedState
}

func createApplication(fromIndex int64, toIndex int64, version string, stateKey string, storedState *storedState, store datastore.DataStore, router Router) (*application, error) {
	out := application{
		fromIndex:   fromIndex,
		toIndex:     toIndex,
		version:     version,
		stateKey:    stateKey,
		storedState: storedState,
		store:       store,
		router:      router,
	}

	return &out, nil
}

// FromBlockIndex returns the from block index
func (app *application) FromBlockIndex() int64 {
	return app.fromIndex
}

// ToBlockIndex returns the to block index
func (app *application) ToBlockIndex() int64 {
	return app.toIndex
}

// GetBlockIndex returns the block index
func (app *application) GetBlockIndex() int64 {
	return app.storedState.State(app.version).Height()
}

// Info returns the application's information
func (app *application) Info(req InfoRequest) InfoResponse {
	version := req.Version()
	state := app.storedState.State(app.version)
	size := state.Size()
	lastBlkHeight := state.Height()
	lastBlkAppHash := state.Hash()
	out := createInfoResponse(size, version, lastBlkHeight, lastBlkAppHash)
	return out
}

// Transact tries to execute a transaction and return its response
func (app *application) Transact(req TransactionRequest) TransactionResponse {
	//execute the transaction:
	resp := app.execTrx(app.store, req)

	//increment the state size:
	app.storedState.State(app.version).Increment()

	//return the response:
	return resp
}

// CheckTransact verifies if a transaction can be executed and return its response
func (app *application) CheckTransact(req TransactionRequest) TransactionResponse {
	//copy the store:
	store := app.store.Copy()

	//execute the transaction without incrementing the state size, then return the response:
	return app.execTrx(store, req)
}

// Commit commits the pending transactions to a block and update the application state, then return its response
func (app *application) Commit() CommitResponse {
	//get the current state:
	st := app.storedState.State(app.version)
	size := st.Size()

	// get the hash from state:
	appHash := st.Hash()

	//if the size is bigger than 0, use the store head hash:
	if size > 0 {
		appHash = app.store.Head().Head().Get()
	}

	//create the updated state:
	newSt := createState(appHash, st.Height()+1, st.Size())
	stJS, stJSErr := json.Marshal(newSt)
	if stJSErr != nil {
		panic(stJSErr)
	}

	//set the updated state:
	amount := app.storedState.Set(app.stateKey, stJS)
	if amount != 1 {
		panic(errors.New("there was a problem while saving the state in the storedState"))
	}

	return createCommitResponse(appHash, newSt.Height())
}

// Query executes a query request on the application
func (app *application) Query(req QueryRequest) QueryResponse {
	outputErrorFn := func(code int, str string) QueryResponse {
		resp, respErr := createEmptyQueryResponse(code, str)
		if respErr != nil {
			panic(respErr)
		}

		return resp
	}

	ptr := req.Pointer()
	from := ptr.From()
	prepHandler := app.router.Route(from, ptr.Path(), Retrieve)
	if prepHandler == nil {
		return outputErrorFn(RouteNotFound, "the router could not find any route for the given query")
	}

	handler := prepHandler.Handler()
	retrieveFunc := handler.Query()
	if retrieveFunc == nil {
		return outputErrorFn(InvalidRoute, "the router found a route for the given query, but its handler had no query func")
	}

	// retrieve the query response:
	queryResponse, queryResponseErr := retrieveFunc(app.store, from, prepHandler.Path(), prepHandler.Params(), req.Signature())
	if queryResponseErr != nil {
		str := fmt.Sprintf("there was an error while executing the query func: %s", queryResponseErr.Error())
		return outputErrorFn(InvalidRequest, str)
	}

	//return the query response:
	return queryResponse
}

func (app *application) execTrx(store datastore.DataStore, req TransactionRequest) TransactionResponse {
	outputErrorFn := func(code int, str string) TransactionResponse {
		trxResp, trxRespErr := createFreeTransactionResponse(code, str)
		if trxRespErr != nil {
			panic(trxRespErr)
		}

		return trxResp
	}

	// if the transaction is a "save-resource-transaction":
	res := req.Resource()
	if res != nil {
		ptr := res.Pointer()
		from := ptr.From()
		prepHandler := app.router.Route(from, ptr.Path(), Save)
		if prepHandler == nil {
			return outputErrorFn(RouteNotFound, "the router could not find any route for the given save transaction")
		}

		handler := prepHandler.Handler()
		saveTrsFunc := handler.SaveTransaction()
		if saveTrsFunc == nil {
			return outputErrorFn(InvalidRoute, "the router found a route for the given transaction, but its handler had no save transaction func")
		}

		trxResponse, trxResponseErr := saveTrsFunc(store, from, prepHandler.Path(), prepHandler.Params(), res.Data(), req.Signature())
		if trxResponseErr != nil {
			str := fmt.Sprintf("there was an error while executing the save transaction func: %s", trxResponseErr.Error())
			return outputErrorFn(InvalidRequest, str)
		}

		return trxResponse
	}

	ptr := req.Pointer()
	from := ptr.From()
	prepHandler := app.router.Route(from, ptr.Path(), Delete)
	if prepHandler == nil {
		return outputErrorFn(RouteNotFound, "the router could not find any route for the given delete transaction")
	}

	handler := prepHandler.Handler()
	delTrsFunc := handler.DeleteTransaction()
	if delTrsFunc == nil {
		return outputErrorFn(InvalidRoute, "the router found a route for the given transaction, but its handler had no delete transaction func")
	}

	trsResponse, trsResponseErr := delTrsFunc(store, from, prepHandler.Path(), prepHandler.Params(), req.Signature())
	if trsResponseErr != nil {
		str := fmt.Sprintf("there was an error while executing the delete transaction func: %s", trsResponseErr.Error())
		return outputErrorFn(InvalidRequest, str)
	}

	return trsResponse
}
