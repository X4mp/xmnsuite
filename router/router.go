package router

type router struct {
	qrRoutes []QueryRoute
	txRoutes []TrxRoute
}

func createRouter(qrRoutes []QueryRoute, txRoutes []TrxRoute) Router {
	out := router{
		qrRoutes: qrRoutes,
		txRoutes: txRoutes,
	}
	return &out
}

// ExecuteQR executes a query route and returns its response:
func (app *router) ExecuteQR(uri string, queryParams map[string]string) Response {
	return nil
}

// ExecuteTR executes a transaction route and returns its response:
func (app *router) ExecuteTR(uri string, queryParams map[string]string, trxData map[string]interface{}) Response {
	return nil
}
