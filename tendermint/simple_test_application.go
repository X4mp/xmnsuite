package tendermint

import (
	"math/rand"

	"github.com/XMNBlockchain/datamint/router"
)

type simpleTestApplication struct {
	data map[string]string
}

func createSimpleTestApplication() router.Router {
	out := simpleTestApplication{
		data: map[string]string{},
	}

	return &out
}

// Start starts the router
func (app *simpleTestApplication) Start() error {
	return nil
}

// Stop stops the router
func (app *simpleTestApplication) Stop() {

}

// Query executes a query route and returns its response:
func (app *simpleTestApplication) Query(req router.Request) router.QueryResponse {
	path := req.Path()
	if el, ok := app.data[path]; ok {
		return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
			IsSuccess:    true,
			IsAuthorized: true,
			IsNFS:        false,
			GazUsed:      int64(rand.Int() % 20),
			Data:         []byte(el),
			Log:          "success",
		})
	}

	return router.SDKFunc.CreateQueryResponse(router.CreateQueryResponseParams{
		IsSuccess:    false,
		IsAuthorized: true,
		IsNFS:        false,
		GazUsed:      0,
		Data:         nil,
		Log:          "not found",
	})
}

// Transact executes a transaction route and returns its response:
func (app *simpleTestApplication) Transact(req router.Request) router.TrxResponse {
	path := req.Path()
	app.data[path] = string(req.Data())
	return router.SDKFunc.CreateTrxResponse(router.CreateTrxResponseParams{
		IsSuccess:    true,
		IsAuthorized: true,
		IsNFS:        false,
		Tags:         map[string][]byte{},
		GazUsed:      int64(rand.Int() % 20),
		Log:          "success",
	})
}

// CheckTrx executes a transaction check route and returns its response:
func (app *simpleTestApplication) CheckTrx(req router.TrxChkRequest) router.TrxChkResponse {
	return router.SDKFunc.CreateTrxChkResponse(router.CreateTrxChkResponseParams{
		CanBeExecuted:   true,
		CanBeAuthorized: true,
		GazWanted:       int64(rand.Int() % 20),
		Log:             "success",
	})
}
