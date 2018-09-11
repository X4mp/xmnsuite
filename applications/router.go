package applications

import (
	"errors"
	"regexp"

	roles "github.com/xmnservices/xmnsuite/roles"
	users "github.com/xmnservices/xmnsuite/users"
	crypto "github.com/tendermint/tendermint/crypto"
)

/*
 * Handler
 */

type handler struct {
	saveTrx SaveTransactionFn
	delTrx  DeleteTransactionFn
	query   QueryFn
}

func createHandler(saveTrx SaveTransactionFn, delTrx DeleteTransactionFn, query QueryFn) (Handler, int, error) {
	if saveTrx != nil {
		return createHandlerWithSaveTransactionFn(saveTrx), Save, nil
	}

	if delTrx != nil {
		return createHandlerWithDeleteTransactionFn(delTrx), Delete, nil
	}

	if query != nil {
		return createHandlerWithQueryFn(query), Retrieve, nil
	}

	return nil, -1, errors.New("one valid func handler is mandatory in order to create an Handler instance")
}

func createHandlerWithSaveTransactionFn(saveTrx SaveTransactionFn) Handler {
	out := handler{
		saveTrx: saveTrx,
		delTrx:  nil,
		query:   nil,
	}

	return &out
}

func createHandlerWithDeleteTransactionFn(delTrx DeleteTransactionFn) Handler {
	out := handler{
		saveTrx: nil,
		delTrx:  delTrx,
		query:   nil,
	}

	return &out
}

func createHandlerWithQueryFn(query QueryFn) Handler {
	out := handler{
		saveTrx: nil,
		delTrx:  nil,
		query:   query,
	}

	return &out
}

// SaveTransaction returns the save transaction func, if any
func (obj *handler) SaveTransaction() SaveTransactionFn {
	return obj.saveTrx
}

// DeleteTransaction returns the delete transaction func, if any
func (obj *handler) DeleteTransaction() DeleteTransactionFn {
	return obj.delTrx
}

// Query returns the query func, if any
func (obj *handler) Query() QueryFn {
	return obj.query
}

// IsWrite returns true if the handler needs a write access, false otherwise
func (obj *handler) IsWrite() bool {
	return obj.delTrx != nil || obj.saveTrx != nil
}

/*
 * PreparedHandler
 */

type preparedHandler struct {
	path   string
	params map[string]string
	handl  Handler
}

func createPreparedHandler(path string, handl Handler) PreparedHandler {
	out := preparedHandler{
		path:   path,
		params: map[string]string{},
		handl:  handl,
	}

	return &out
}

func createPreparedHandlerWithParams(path string, params map[string]string, handl Handler) PreparedHandler {
	out := preparedHandler{
		path:   path,
		params: params,
		handl:  handl,
	}

	return &out
}

// Path returns the path
func (obj *preparedHandler) Path() string {
	return obj.path
}

// Params returns the params
func (obj *preparedHandler) Params() map[string]string {
	return obj.params
}

// Handler returns the handler
func (obj *preparedHandler) Handler() Handler {
	return obj.handl
}

/*
 * Route
 */

type route struct {
	rols          roles.Roles
	usrs          users.Users
	pattern       *regexp.Regexp
	variableNames []string
	handl         Handler
	roleKey       string
}

func createRoute(roleKey string, rols roles.Roles, usrs users.Users, patternAsString string, handl Handler) (Route, error) {
	//create the pattern:
	pattern, variableNames, patternErr := fromURLPatternToRegex(patternAsString)
	if patternErr != nil {
		return nil, patternErr
	}

	//create the role:

	out := route{
		rols:          rols,
		usrs:          usrs,
		pattern:       pattern,
		variableNames: variableNames,
		handl:         handl,
		roleKey:       roleKey,
	}

	return &out, nil
}

// Matches returns true if the route matches the regex, false otherwise
func (obj *route) Matches(from crypto.PubKey, path string) bool {

	//if the route needs write access:
	if obj.handl.IsWrite() {
		writeAccessKeys := obj.rols.HasWriteAccess(obj.roleKey, path)
		if len(writeAccessKeys) <= 0 {
			return false
		}

		return true
	}

	foundStr := obj.pattern.FindString(path)
	if foundStr != path {
		return false
	}

	return true
}

// Handler returns the handler
func (obj *route) Handler(from crypto.PubKey, path string) PreparedHandler {
	if !obj.Matches(from, path) {
		return nil
	}

	valuesWithURL := obj.pattern.FindStringSubmatch(path)
	values := valuesWithURL[1:]
	if len(values) != len(obj.variableNames) {
		return nil
	}

	params := map[string]string{}
	for index, oneVariableName := range obj.variableNames {
		params[oneVariableName] = values[index]
	}

	out := createPreparedHandlerWithParams(path, params, obj.handl)
	return out
}

/*
 * Router
 */
type router struct {
	rtes map[int][]Route
}

func createRouter(rtes map[int][]Route) Router {
	out := router{
		rtes: rtes,
	}

	return &out
}

// Route route a request
func (obj *router) Route(from crypto.PubKey, path string, method int) PreparedHandler {
	if rtes, ok := obj.rtes[method]; ok {
		for _, oneRte := range rtes {
			handl := oneRte.Handler(from, path)
			if handl != nil {
				return handl
			}
		}

		return nil
	}

	return nil
}
