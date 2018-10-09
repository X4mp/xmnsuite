package tendermint

import (
	"fmt"

	types "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	routers "github.com/xmnservices/xmnsuite/routers"
)

/*
 * ABCI Application
 */

type abciApplication struct {
	types.BaseApplication
	apps      applications.Applications
	blkHeight int64
}

func createABCIApplication(apps applications.Applications) (*abciApplication, error) {
	out := abciApplication{
		apps:      apps,
		blkHeight: apps.RetrieveBlockIndex(),
	}

	return &out, nil
}

// Info outputs information related to the abciApplication state
func (app *abciApplication) Info(req types.RequestInfo) types.ResponseInfo {
	// retrieve the app:
	curApp, curAppErr := app.apps.RetrieveByBlockIndex(app.blkHeight)
	if curAppErr != nil {
		panic(curAppErr)
	}

	// execute the request on the application:
	resp := curApp.Info(applications.SDKFunc.CreateInfoRequest(applications.CreateInfoRequestParams{
		Version: req.GetVersion(),
	}))

	out := struct {
		Size int64 `json:"size"`
	}{
		Size: resp.State().Size(),
	}

	js, jsErr := cdc.MarshalJSON(out)
	if jsErr != nil {
		panic(js)
	}

	fmt.Printf("Info last height: %d, last AppHash: %X\n", app.blkHeight, resp.State().Hash())

	if resp.State().Size() > 0 {
		return types.ResponseInfo{
			Data:             string(js),
			Version:          resp.Version(),
			LastBlockHeight:  resp.State().Height(),
			LastBlockAppHash: resp.State().Hash(),
		}
	}

	return types.ResponseInfo{
		Data:            string(js),
		Version:         resp.Version(),
		LastBlockHeight: resp.State().Height(),
	}
}

// DeliverTx delivers a transaction to the abciApplication
func (app *abciApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	// retrieve the app:
	curApp, curAppErr := app.apps.RetrieveByBlockIndex(app.blkHeight)
	if curAppErr != nil {
		panic(curAppErr)
	}

	// execute the transaction on the application:
	resp := curApp.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
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
	// retrieve the app:
	curApp, curAppErr := app.apps.RetrieveByBlockIndex(app.blkHeight)
	if curAppErr != nil {
		panic(curAppErr)
	}

	// execute the transaction on the application:
	resp := curApp.CheckTransact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
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
	// retrieve the app:
	curApp, curAppErr := app.apps.RetrieveByBlockIndex(app.blkHeight)
	if curAppErr != nil {
		panic(curAppErr)
	}

	//execute the commit on the application:
	resp := curApp.Commit()

	// fetch the data from the response:
	appHash := resp.AppHash()

	// update the block height:
	app.blkHeight = resp.BlockHeight()

	fmt.Printf("Commit height: %d, AppHash: %X\n", app.blkHeight, appHash)

	// return the value:
	return types.ResponseCommit{Data: appHash}
}

// Query executes a query on the abciApplication
func (app *abciApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	if !reqQuery.GetProve() {
		return types.ResponseQuery{
			Code: uint32(routers.InvalidRequest),
			Log:  "the query cannot be trusted",
		}
	}

	// retrieve the app:
	curApp, curAppErr := app.apps.RetrieveByBlockIndex(app.blkHeight)
	if curAppErr != nil {
		panic(curAppErr)
	}

	//execute the query on the application:
	resp := curApp.Query(routers.SDKFunc.CreateQueryRequest(routers.CreateQueryRequestParams{
		JSData: reqQuery.GetData(),
	}))

	//fetch the data from the response:
	code := resp.Code()
	key := resp.Key()
	value := resp.Value()
	log := resp.Log()

	// return the value:
	out := types.ResponseQuery{
		Code: uint32(code),
		Log:  log,
	}

	if key != "" {
		out.Key = []byte(key)
	}

	if value != nil {
		out.Value = value
	}

	return out
}
