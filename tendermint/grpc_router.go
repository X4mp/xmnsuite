package tendermint

import (
	"time"

	"github.com/XMNBlockchain/datamint/router"
	client "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

type grpcRouter struct {
	cl client.Client
}

func createGRPCRouter(ipAddress string) (router.Router, error) {

	healthy := func(ipAddress string) bool {
		//connect to rpc:
		rpc := rpcclient.NewJSONRPCClient(ipAddress)
		ctypes.RegisterAmino(rpc.Codec())

		result := new(ctypes.ResultHealth)
		_, outErr := rpc.Call("health", nil, result)
		if outErr == nil {
			return true
		}

		return false
	}

	//wait until healthy:
	for {
		if healthy(ipAddress) {
			break
		}

		time.Sleep(time.Second)
	}

	cl, clErr := client.NewClient(ipAddress, "grpc", true)
	if clErr != nil {
		return nil, clErr
	}

	out := grpcRouter{
		cl: cl,
	}

	return &out, nil
}

// Query executes a query route and returns its response:
func (app *grpcRouter) Query(req router.Request) router.QueryResponse {
	js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		panic(jsErr)
	}

	reqQuery := types.RequestQuery{
		Data: js,
	}

	clResp, clRespErr := app.cl.QuerySync(reqQuery)
	if clRespErr != nil {
		return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
			IsSuccess: false,
			Log:       clResp.GetLog(),
		})
	}

	return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
		JSData: clResp.GetValue(),
	})
}

// Transact executes a transaction route and returns its response:
func (app *grpcRouter) Transact(req router.Request) router.TrxResponse {
	js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		panic(jsErr)
	}

	clResp, clRespErr := app.cl.DeliverTxSync(js)
	if clRespErr != nil {
		return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
			IsSuccess: false,
			Log:       clResp.GetLog(),
		})
	}

	return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
		JSData: clResp.GetData(),
	})
}

// CheckTrx executes a transaction check route and returns its response:
func (app *grpcRouter) CheckTrx(req router.TrxChkRequest) router.TrxChkResponse {
	js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		panic(jsErr)
	}

	clResp, clRespErr := app.cl.CheckTxSync(js)
	if clRespErr != nil {
		return router.SDKFunc.CreateTrxChkResponse(router.CreateTrxChkResponseParams{
			CanBeExecuted: false,
			Log:           clResp.GetLog(),
		})
	}

	return router.SDKFunc.CreateTrxChkResponse(router.CreateTrxChkResponseParams{
		JSData: clResp.GetData(),
	})
}
