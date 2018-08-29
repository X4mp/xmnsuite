package tendermint

import (
	"encoding/binary"
	"encoding/json"

	applications "github.com/XMNBlockchain/datamint/applications"
	hashtree "github.com/XMNBlockchain/datamint/hashtree"
	router "github.com/XMNBlockchain/datamint/router"
	code "github.com/tendermint/tendermint/abci/example/code"
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
	//execute the transaction:
	resp := app.rt.Transact(router.SDKFunc.CreateRequest(router.CreateRequestParams{
		JSData: tx,
	}))

	//if the transaction is un-authorized:
	if !resp.IsAuthorized() {
		return types.ResponseDeliverTx{Code: code.CodeTypeUnauthorized, Log: resp.Log(), GasUsed: resp.GazUsed()}
	}

	//if the transaction was not successful, flag the encoding error:
	if !resp.IsSuccess() {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError, Log: resp.Log(), GasUsed: resp.GazUsed()}
	}

	//create the tags:
	tags := resp.Tags()
	tagPairs := []cmn.KVPair{}
	for key, value := range tags {
		tagPairs = append(tagPairs, cmn.KVPair{
			Key:   []byte(key),
			Value: value,
		})
	}

	//the transaction was successful:
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Log: resp.Log(), GasUsed: resp.GazUsed(), Tags: tagPairs}
}

// CheckTx verifies that a transaction is valid before it gets executed
func (app *abciApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	//execute the transaction:
	/*resp := app.rt.CheckTrx(router.SDKFunc.CreateRequest(router.CreateRequestParams{
		JSData: tx,
	}))

	if !resp.CanBeAuthorized() {
		return types.ResponseCheckTx{Code: code.CodeTypeUnauthorized, Log: resp.Log(), GasWanted: resp.GazWanted()}
	}

	if !resp.CanBeExecuted() {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError, Log: resp.Log(), GasWanted: resp.GazWanted()}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, Log: resp.Log(), GasWanted: resp.GazWanted()}*/
	return nil
}

// Commit commits the blockchain
func (app *abciApplication) Commit() types.ResponseCommit {

	//get the current state:
	st := app.state.GetState()
	size := app.state.GetState().GetSize()

	//generate an app hash:
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, size)

	//if the size is bigger than 0, add the store head:
	if size > 0 {
		head := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
			Blocks: [][]byte{
				appHash,
				app.store.Head().Head().Get(),
			},
		})

		appHash = head.Head().Get()
	}

	//create the updated state:
	newSt := createState(appHash, st.GetKeynamePrefix(), st.GetHeight()+1, st.GetSize())
	stJS, stJSErr := json.Marshal(newSt)
	if stJSErr != nil {
		panic(stJSErr)
	}

	//set the updated state:
	app.state.Set(app.stateKey, stJS)

	//return the response:
	return types.ResponseCommit{Data: appHash}
}

// Query executes a query on the abciApplication
func (app *abciApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	resp := app.rt.Query(router.SDKFunc.CreateRequest(router.CreateRequestParams{
		JSData: reqQuery.GetData(),
	}))

	js, jsErr := cdc.MarshalJSON(resp)
	if jsErr != nil {
		panic(jsErr)
	}

	return types.ResponseQuery{
		Value: js,
		Log:   "success",
	}
}
