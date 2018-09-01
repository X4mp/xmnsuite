package tendermint

import (
	applications "github.com/XMNBlockchain/datamint/applications"
	types "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

/*
 * ABCI Application
 */

type abciApplication struct {
	types.BaseApplication
	app applications.Application
}

func createABCIApplication(app applications.Application) (*abciApplication, error) {
	out := abciApplication{
		app: app,
	}

	return &out, nil
}

// Info outputs information related to the abciApplication state
func (app *abciApplication) Info(req types.RequestInfo) types.ResponseInfo {

	// execute the request on the application:
	resp := app.app.Info(applications.SDKFunc.CreateInfoRequest(applications.CreateInfoRequestParams{
		Version: req.GetVersion(),
	}))

	out := struct {
		Size int64 `json:"size"`
	}{
		Size: resp.Size(),
	}

	js, jsErr := cdc.MarshalJSON(out)
	if jsErr != nil {
		panic(js)
	}

	return types.ResponseInfo{
		Data:             string(js),
		Version:          resp.Version(),
		LastBlockHeight:  resp.LastBlockHeight(),
		LastBlockAppHash: resp.LastBlockAppHash(),
	}
}

// DeliverTx delivers a transaction to the abciApplication
func (app *abciApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	// execute the transaction on the application:
	resp := app.app.Transact(applications.SDKFunc.CreateTransactionRequest(applications.CreateTransactionRequestParams{
		JSData: tx,
	}))

	// fetch the data from response:
	code := resp.Code()
	gazUsed := resp.GazUsed()
	log := resp.Log()
	inputTags := resp.Tags()

	//create the tags:
	tagPairs := []cmn.KVPair{}
	for key, value := range inputTags {
		tagPairs = append(tagPairs, cmn.KVPair{
			Key:   []byte(key),
			Value: value,
		})
	}

	//return the value:
	return types.ResponseDeliverTx{Code: uint32(code), Log: log, GasUsed: gazUsed, Tags: tagPairs}
}

// CheckTx verifies that a transaction is valid before it gets executed
func (app *abciApplication) CheckTx(tx []byte) types.ResponseCheckTx {

	// execute the transaction on the application:
	resp := app.app.CheckTransact(applications.SDKFunc.CreateTransactionRequest(applications.CreateTransactionRequestParams{
		JSData: tx,
	}))

	// fetch the data from response:
	code := resp.Code()
	gazWanted := resp.GazUsed()
	log := resp.Log()
	inputTags := resp.Tags()

	//create the tags:
	tagPairs := []cmn.KVPair{}
	for key, value := range inputTags {
		tagPairs = append(tagPairs, cmn.KVPair{
			Key:   []byte(key),
			Value: value,
		})
	}

	//return the value:
	return types.ResponseCheckTx{Code: uint32(code), Log: log, GasWanted: gazWanted, Tags: tagPairs}
}

// Commit commits the blockchain
func (app *abciApplication) Commit() types.ResponseCommit {
	//execute the commit on the application:
	resp := app.app.Commit()

	// fetch the data from the response:
	appHash := resp.AppHash()

	// return the value:
	return types.ResponseCommit{Data: appHash}
}

// Query executes a query on the abciApplication
func (app *abciApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {

	if !reqQuery.GetProve() {
		return types.ResponseQuery{
			Code: uint32(applications.InvalidRequest),
			Log:  "the query must be proved",
		}
	}

	//execute the query on the application:
	resp := app.app.Query(applications.SDKFunc.CreateQueryRequest(applications.CreateQueryRequestParams{
		JSData: reqQuery.GetData(),
	}))

	//fetch the data from the response:
	code := resp.Code()
	key := resp.Key()
	value := resp.Value()
	log := resp.Log()

	// return the value:
	return types.ResponseQuery{
		Code:  uint32(code),
		Log:   log,
		Key:   []byte(key),
		Value: value,
	}
}
