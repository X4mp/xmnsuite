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

// Start starts the router
func (app *router) Start() error {
	return nil
}

// Stop stops the router
func (app *router) Stop() {

}

// Query executes a query route and returns its response:
func (app *router) Query(req Request) QueryResponse {
	for _, oneRoute := range app.qrRoutes {
		if oneRoute.Matches(req) {
			fn := oneRoute.Handler()
			return fn(req)
		}
	}

	return nil
}

// Transact executes a transaction route and returns its response:
func (app *router) Transact(req Request) TrxResponse {
	for _, oneRoute := range app.txRoutes {
		if oneRoute.Matches(req) {
			fn := oneRoute.Handler()
			return fn(req)
		}
	}

	return nil
}
