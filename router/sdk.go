package router

import "regexp"

// QueryHandlerFn represents a query handler func
type QueryHandlerFn func(uriParams map[string]string, queryParams map[string]string) Response

// TrxHandlerFn represents a transaction handler func
type TrxHandlerFn func(uriParams map[string]string, queryParams map[string]string, trxData map[string]interface{}) Response

// Header represents the header of a response
type Header interface {
	StatusCode() int
	Lines() map[string]string
}

// Response represents an handler response
type Response interface {
	Header() Header
	Body() []byte
}

// QueryRoute represents a query route
type QueryRoute interface {
	Pattern() *regexp.Regexp
	Handler() QueryHandlerFn
}

// TrxRoute represents a transaction route
type TrxRoute interface {
	Pattern() *regexp.Regexp
	Handler() TrxHandlerFn
}

// Router represents the router
type Router interface {
	RegisterQR(routes ...QueryRoute) int
	RegisterTR(routes ...TrxRoute) int
	ExecuteQR(uri string, queryParams map[string]string) Response
	ExecuteTR(uri string, queryParams map[string]string, trxData map[string]interface{}) Response
}
