package router

type router struct {
	qrRoutes    []QueryRoute
	txChkRoutes []TrxChkRoute
	txRoutes    []TrxRoute
}

func createRouter(qrRoutes []QueryRoute, txChkRoutes []TrxChkRoute, txRoutes []TrxRoute) Router {
	out := router{
		qrRoutes:    qrRoutes,
		txChkRoutes: txChkRoutes,
		txRoutes:    txRoutes,
	}
	return &out
}

// Query executes a query route and returns its response:
func (app *router) Query(req Request) QueryResponse {
	return nil
}

// Transact executes a transaction route and returns its response:
func (app *router) Transact(req Request) TrxResponse {
	return nil
}

// CheckTrx executes a transaction check route and returns its response:
func (app *router) CheckTrx(req TrxChkRequest) TrxChkResponse {
	return nil
}
