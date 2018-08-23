package tendermint

import (
	"net"
	"time"

	"github.com/XMNBlockchain/datamint/router"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/abci/types"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	cmn "github.com/tendermint/tendermint/libs/common"
)

/*
 * GRPC Codec
 */

type grpcCodec struct {
	cdc *amino.Codec
}

func createGRPCCodec(cdc *amino.Codec) *grpcCodec {
	out := grpcCodec{
		cdc: cdc,
	}

	return &out
}

// Marshal marshals an object
func (obj *grpcCodec) Marshal(v interface{}) ([]byte, error) {
	return cdc.MarshalJSON(v)
}

// Unmarshal unmarshals an object
func (obj *grpcCodec) Unmarshal(data []byte, v interface{}) error {
	return cdc.UnmarshalJSON(data, v)
}

// String returns the name of the Codec implementation.
func (obj *grpcCodec) String() string {
	return "grpc-codec"
}

/*
 * GRPC Router
 */

type grpcRouter struct {
	cl        types.ABCIApplicationClient
	ipAddress string
	conn      *grpc.ClientConn
}

func createGRPCRouter(ipAddress string) (router.Router, error) {
	out := grpcRouter{
		cl:        nil,
		conn:      nil,
		ipAddress: ipAddress,
	}

	return &out, nil
}

// Start starts the router
func (app *grpcRouter) Start() error {

	//setup the codec in GRPC:
	cdcDialOpt := grpc.WithCodec(createGRPCCodec(cdc))

	// Connect to the socket
	conn, err := grpc.Dial(app.ipAddress, grpc.WithInsecure(), cdcDialOpt, grpc.WithDialer(func(addr string, dur time.Duration) (net.Conn, error) {
		return cmn.Connect(addr)
	}))

	if err != nil {
		return err
	}

	app.cl = types.NewABCIApplicationClient(conn)
	app.conn = conn
	return nil
}

// Stop stops the router
func (app *grpcRouter) Stop() {
	err := app.conn.Close()
	if err != nil {
		panic(err)
	}
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

	clResp, clRespErr := app.cl.Query(context.Background(), &reqQuery)
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

	tx := types.RequestDeliverTx{
		Tx: js,
	}

	clResp, clRespErr := app.cl.DeliverTx(context.Background(), &tx)
	if clRespErr != nil {
		return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
			IsSuccess: false,
			Log:       clRespErr.Error(),
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

	tx := types.RequestCheckTx{
		Tx: js,
	}

	clResp, clRespErr := app.cl.CheckTx(context.Background(), &tx)
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
