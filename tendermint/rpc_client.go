package tendermint

import (
	"fmt"

	applications "github.com/XMNBlockchain/datamint/applications"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

/*
 * RPC Router
 */

type rpcClient struct {
	ipAddress string
	cl        *rpcclient.JSONRPCClient
}

func createRPCClient(ipAddress string) (applications.Client, error) {
	out := rpcClient{
		ipAddress: ipAddress,
		cl:        nil,
	}

	return &out, nil
}

// Start starts the router
func (app *rpcClient) Start() error {
	//create the client set the codec:
	client := rpcclient.NewJSONRPCClient(app.ipAddress)
	client.SetCodec(cdc)

	//keep the instance:
	app.cl = client
	return nil
}

// Query executes a query and returns its response:
func (app *rpcClient) Query(req applications.QueryRequest) (applications.QueryResponse, error) {
	js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		return nil, jsErr
	}

	params := map[string]interface{}{
		"path":    req.Pointer().Path(),
		"data":    fmt.Sprintf("%X", js),
		"height":  0,
		"trusted": false,
	}

	result := new(ctypes.ResultABCIQuery)
	_, outErr := app.cl.Call("abci_query", params, result)
	if outErr != nil {
		return nil, outErr
	}

	return applications.SDKFunc.CreateQueryResponse(applications.CreateQueryResponseParams{
		Code:  int(result.Response.GetCode()),
		Log:   result.Response.GetLog(),
		Key:   string(result.Response.GetKey()),
		Value: result.Response.GetValue(),
	}), nil
}

// Transact executes a transaction and returns its response:
func (app *rpcClient) Transact(req applications.TransactionRequest) (applications.ClientTransactionResponse, error) {
	reqJS, reqJSErr := cdc.MarshalJSON(req)
	if reqJSErr != nil {
		return nil, reqJSErr
	}

	result := new(ctypes.ResultBroadcastTxCommit)
	_, err := app.cl.Call("broadcast_tx_commit", map[string]interface{}{"tx": reqJS}, result)
	if err != nil {
		return nil, err
	}

	// retrieve the transaction data:
	code := result.DeliverTx.GetCode()
	log := result.DeliverTx.GetLog()
	gazUsed := result.DeliverTx.GetGasUsed()

	tags := map[string][]byte{}
	pairs := result.DeliverTx.GetTags()
	for _, onePair := range pairs {
		tags[string(onePair.GetKey())] = onePair.GetValue()
	}

	// create the transaction response:
	trsResponse := applications.SDKFunc.CreateTransactionResponse(applications.CreateTransactionResponseParams{
		Code:    int(code),
		Log:     log,
		GazUsed: gazUsed,
		Tags:    tags,
	})

	// retrieve the check data:
	chkCode := result.CheckTx.GetCode()
	chkLog := result.CheckTx.GetLog()
	chkGazUsed := result.CheckTx.GetGasUsed()

	chkTags := map[string][]byte{}
	chkPairs := result.CheckTx.GetTags()
	for _, onePair := range chkPairs {
		chkTags[string(onePair.GetKey())] = onePair.GetValue()
	}

	// retrieve the commit data:
	chkResponse := applications.SDKFunc.CreateTransactionResponse(applications.CreateTransactionResponseParams{
		Code:    int(chkCode),
		Log:     chkLog,
		GazUsed: chkGazUsed,
		Tags:    chkTags,
	})

	return applications.SDKFunc.CreateClientTransactionResponse(applications.CreateClientTransactionResponseParams{
		Chk:    chkResponse,
		Trx:    trsResponse,
		Height: result.Height,
		Hash:   result.Hash,
	}), nil
}
