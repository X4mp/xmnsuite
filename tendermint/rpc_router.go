package tendermint

import (
	"github.com/XMNBlockchain/datamint/router"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

/*
 * RPC Router
 */

type rpcRouter struct {
	ipAddress string
	cl        *rpcclient.JSONRPCClient
}

func createRPCRouter(ipAddress string) (router.Router, error) {
	out := rpcRouter{
		ipAddress: ipAddress,
		cl:        nil,
	}

	return &out, nil
}

// Start starts the router
func (app *rpcRouter) Start() error {
	//create the client set the codec:
	client := rpcclient.NewJSONRPCClient(app.ipAddress)
	client.SetCodec(cdc)

	//keep the instance:
	app.cl = client
	return nil
}

// Stop stops the router
func (app *rpcRouter) Stop() {
	app.cl = nil
}

// Query executes a query route and returns its response:
func (app *rpcRouter) Query(req router.Request) router.QueryResponse {
	/*js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		panic(jsErr)
	}

	reqQuery := types.RequestQuery{
		Data: js,
	}

	clResp, clRespErr := app.cl.Query(&reqQuery)
	if clRespErr != nil {
		return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
			IsSuccess: false,
			Log:       clResp.GetLog(),
		})
	}

	return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
		JSData: clResp.GetValue(),
	})*/
	return nil
}

// Transact executes a transaction route and returns its response:
func (app *rpcRouter) Transact(req router.Request) router.TrxResponse {
	reqJS, reqJSErr := cdc.MarshalJSON(req)
	if reqJSErr != nil {
		return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
			IsSuccess: false,
			Log:       reqJSErr.Error(),
		})
	}

	result := new(ctypes.ResultBroadcastTxCommit)
	_, err := app.cl.Call("broadcast_tx_commit", map[string]interface{}{"tx": reqJS}, result)
	if err != nil {
		return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
			IsSuccess: false,
			Log:       err.Error(),
		})
	}

	tags := map[string][]byte{}
	pairs := result.DeliverTx.GetTags()
	for _, onePair := range pairs {
		key := string(onePair.GetKey())
		tags[key] = onePair.GetValue()
	}

	return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
		IsSuccess:    result.DeliverTx.IsOK(),
		IsAuthorized: true,
		IsNFS:        true,
		Tags:         tags,
		GazUsed:      result.DeliverTx.GetGasUsed(),
		Log:          result.DeliverTx.GetLog(),
	})
}
