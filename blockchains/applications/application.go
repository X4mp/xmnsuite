package applications

import (
	"fmt"
	"log"

	datastore "github.com/xmnservices/xmnsuite/datastore"
	routers "github.com/xmnservices/xmnsuite/routers"
)

/*
 * Application
 */

type application struct {
	fromIndex int64
	toIndex   int64
	version   string
	stateKey  string
	router    routers.Router
	db        Database
}

func createApplication(
	fromIndex int64,
	toIndex int64,
	version string,
	db Database,
	router routers.Router,
) (*application, error) {
	out := application{
		fromIndex: fromIndex,
		toIndex:   toIndex,
		version:   version,
		db:        db,
		router:    router,
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
	return app.db.State(app.version).Height()
}

// Info returns the application's information
func (app *application) Info(req InfoRequest) InfoResponse {
	version := req.Version()
	state := app.db.State(app.version)
	out := createInfoResponse(version, state)
	return out
}

// Transact tries to execute a transaction and return its response
func (app *application) Transact(req routers.TransactionRequest) routers.TransactionResponse {
	//execute the transaction:
	resp := app.execTrx(app.db.DataStore().DataStore(), req)

	//increment the state size:
	app.db.State(app.version).Increment()

	//return the response:
	return resp
}

// CheckTransact verifies if a transaction can be executed and return its response
func (app *application) CheckTransact(req routers.TransactionRequest) routers.TransactionResponse {
	//copy the store:
	store := app.db.DataStore().DataStore().Copy()

	//execute the transaction without incrementing the state size, then return the response:
	return app.execTrx(store, req)
}

// Commit commits the pending transactions to a block and update the application state, then return its response
func (app *application) Commit() CommitResponse {
	// current state:
	curSt := app.db.State(app.version)

	// update the state:
	st, stErr := app.db.Update(app.version)
	if stErr != nil {
		panic(stErr)
	}

	// response:
	return createCommitResponse(curSt.Hash(), st.Hash(), st.Height())
}

// Query executes a query request on the application
func (app *application) Query(req routers.QueryRequest) routers.QueryResponse {

	defer func() {
		if r := recover(); r != nil {
			log.Println("There was an error while executing the query:", r)
		}
	}()

	outputErrorFn := func(code int, str string) routers.QueryResponse {
		resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code: code,
			Log:  str,
		})

		return resp
	}

	ptr := req.Pointer()
	from := ptr.From()
	prepHandler := app.router.Route(from, ptr.Path(), routers.Retrieve)
	if prepHandler == nil {
		return outputErrorFn(routers.RouteNotFound, "the router could not find any route for the given query")
	}

	handler := prepHandler.Handler()
	retrieveFunc := handler.Query()
	if retrieveFunc == nil {
		return outputErrorFn(routers.InvalidRoute, "the router found a route for the given query, but its handler had no query func")
	}

	// retrieve the query response:
	queryResponse, queryResponseErr := retrieveFunc(app.db.DataStore().DataStore(), from, prepHandler.Path(), prepHandler.Params(), req.Signature())
	if queryResponseErr != nil {
		str := fmt.Sprintf("there was an error while executing the query func: %s", queryResponseErr.Error())
		return outputErrorFn(routers.InvalidRequest, str)
	}

	//return the query response:
	return queryResponse
}

func (app *application) execTrx(store datastore.DataStore, req routers.TransactionRequest) routers.TransactionResponse {

	defer func() {
		if r := recover(); r != nil {
			log.Println("There was an error while executing the transaction:", r)
		}
	}()

	outputErrorFn := func(code int, str string) routers.TransactionResponse {
		trxResp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
			Code: code,
			Log:  str,
		})

		return trxResp
	}

	// if the transaction is a "save-resource-transaction":
	res := req.Resource()
	if res != nil {
		ptr := res.Pointer()
		from := ptr.From()
		prepHandler := app.router.Route(from, ptr.Path(), routers.Save)
		if prepHandler == nil {
			return outputErrorFn(routers.RouteNotFound, "the router could not find any route for the given save transaction")
		}

		handler := prepHandler.Handler()
		saveTrsFunc := handler.SaveTransaction()
		if saveTrsFunc == nil {
			return outputErrorFn(routers.InvalidRoute, "the router found a route for the given transaction, but its handler had no save transaction func")
		}

		trxResponse, trxResponseErr := saveTrsFunc(store, from, prepHandler.Path(), prepHandler.Params(), res.Data(), req.Signature())
		if trxResponseErr != nil {
			str := fmt.Sprintf("there was an error while executing the save transaction func: %s", trxResponseErr.Error())
			return outputErrorFn(routers.InvalidRequest, str)
		}

		return trxResponse
	}

	ptr := req.Pointer()
	from := ptr.From()
	prepHandler := app.router.Route(from, ptr.Path(), routers.Delete)
	if prepHandler == nil {
		return outputErrorFn(routers.RouteNotFound, "the router could not find any route for the given delete transaction")
	}

	handler := prepHandler.Handler()
	delTrsFunc := handler.DeleteTransaction()
	if delTrsFunc == nil {
		return outputErrorFn(routers.InvalidRoute, "the router found a route for the given transaction, but its handler had no delete transaction func")
	}

	trsResponse, trsResponseErr := delTrsFunc(store, from, prepHandler.Path(), prepHandler.Params(), req.Signature())
	if trsResponseErr != nil {
		str := fmt.Sprintf("there was an error while executing the delete transaction func: %s", trsResponseErr.Error())
		return outputErrorFn(routers.InvalidRequest, str)
	}

	return trsResponse
}
