package applications

import (
	"fmt"
	"reflect"
	"testing"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	roles "github.com/xmnservices/xmnsuite/roles"
	users "github.com/xmnservices/xmnsuite/users"
)

func TestCreateHandler_withSaveTrxFn_Success(t *testing.T) {
	//variables:
	fn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (TransactionResponse, error) {
		return nil, nil
	}

	//execute:
	handler := createHandlerWithSaveTransactionFn(fn)
	saveTrx := handler.SaveTransaction()
	delTrx := handler.DeleteTransaction()
	query := handler.Query()
	isWrite := handler.IsWrite()

	if saveTrx == nil {
		t.Errorf("the returned SaveTransaction func was expected to be a func, nil returned")
		return
	}

	if delTrx != nil {
		t.Errorf("the returned DeleteTransaction func was expected to be nil, func returned")
		return
	}

	if query != nil {
		t.Errorf("the returned Query func was expected to be nil, func returned")
		return
	}

	if !isWrite {
		t.Errorf("the handler was expected to need write acess")
		return
	}
}

func TestCreateHandler_withDeleteTrxFn_Success(t *testing.T) {
	//variables:
	fn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (TransactionResponse, error) {
		return nil, nil
	}

	//execute:
	handler := createHandlerWithDeleteTransactionFn(fn)
	saveTrx := handler.SaveTransaction()
	delTrx := handler.DeleteTransaction()
	query := handler.Query()
	isWrite := handler.IsWrite()

	if saveTrx != nil {
		t.Errorf("the returned SaveTransaction func was expected to be nil, func returned")
		return
	}

	if delTrx == nil {
		t.Errorf("the returned DeleteTranaction func was expected to be a func, nil returned")
		return
	}

	if query != nil {
		t.Errorf("the returned Query func was expected to be nil, func returned")
		return
	}

	if !isWrite {
		t.Errorf("the handler was expected to need write acess")
		return
	}
}

func TestCreateHandler_withQueryFn_Success(t *testing.T) {
	//variables:
	fn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		return nil, nil
	}

	//execute:
	handler := createHandlerWithQueryFn(fn)
	saveTrx := handler.SaveTransaction()
	delTrx := handler.DeleteTransaction()
	query := handler.Query()
	isWrite := handler.IsWrite()

	if saveTrx != nil {
		t.Errorf("the returned SaveTransaction func was expected to be nil, func returned")
		return
	}

	if delTrx != nil {
		t.Errorf("the returned DeleteTransaction func was expected to be nil, func returned")
		return
	}

	if query == nil {
		t.Errorf("the returned Query func was expected to be a func, nil returned")
		return
	}

	if isWrite {
		t.Errorf("the handler was expected to NOT need write acess")
		return
	}
}

func TestCreatePreparedHandler_Success(t *testing.T) {
	//variables:
	queryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		return nil, nil
	}

	path := "/this/is/a/path"
	handler := createHandlerWithQueryFn(queryFn)

	//execute:
	preparedHandler := createPreparedHandler(path, handler)
	retPath := preparedHandler.Path()
	retHandler := preparedHandler.Handler()
	retParams := preparedHandler.Params()

	if !reflect.DeepEqual(path, retPath) {
		t.Errorf("the returned path is invalid.  Expected: %s, Returned: %s", path, retPath)
		return
	}

	if !reflect.DeepEqual(handler, retHandler) {
		t.Errorf("the returned handler is invalid.")
		return
	}

	if !reflect.DeepEqual(map[string]string{}, retParams) {
		t.Errorf("the returned params are invalid.")
		return
	}
}

func TestCreatePreparedHandler_withParams_Success(t *testing.T) {
	//variables:
	queryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		return nil, nil
	}

	path := "/this/is/a/path"
	handler := createHandlerWithQueryFn(queryFn)
	params := map[string]string{
		"some": "params",
	}

	//execute:
	preparedHandler := createPreparedHandlerWithParams(path, params, handler)
	retPath := preparedHandler.Path()
	retHandler := preparedHandler.Handler()
	retParams := preparedHandler.Params()

	if !reflect.DeepEqual(path, retPath) {
		t.Errorf("the returned path is invalid.  Expected: %s, Returned: %s", path, retPath)
		return
	}

	if !reflect.DeepEqual(handler, retHandler) {
		t.Errorf("the returned handler is invalid.")
		return
	}

	if !reflect.DeepEqual(params, retParams) {
		t.Errorf("the returned params are invalid.")
		return
	}
}

func TestCreateRoute_withReadRoute_matches_Success(t *testing.T) {
	//variables:
	queryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		return nil, nil
	}

	rols := roles.SDKFunc.Create()
	usrs := users.SDKFunc.Create()
	roleKey := "video-update-role-01"
	patternAsString := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	handler := createHandlerWithQueryFn(queryFn)
	path := fmt.Sprintf("/videos/%s", "6adbfdfc-bb7d-4236-96d6-96d1688a2441")

	//execute:
	route, routeErr := createRoute(roleKey, rols, usrs, patternAsString, handler)
	if routeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", routeErr.Error())
		return
	}

	matches := route.Matches(nil, path)
	if !matches {
		t.Errorf("the route was expected to match")
		return
	}

	retHandler := route.Handler(nil, path)
	if retHandler == nil {
		t.Errorf("the returned handler was expected to be valid, nil returned")
		return
	}
}

func TestCreateRoute_withReadRoute_doesNotMatch_Success(t *testing.T) {
	//variables:
	queryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		return nil, nil
	}

	rols := roles.SDKFunc.Create()
	usrs := users.SDKFunc.Create()
	roleKey := "video-update-role-01"
	patternAsString := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	handler := createHandlerWithQueryFn(queryFn)
	path := fmt.Sprintf("/nomatch/%s", "6adbfdfc-bb7d-4236-96d6-96d1688a2441")

	//execute:
	route, routeErr := createRoute(roleKey, rols, usrs, patternAsString, handler)
	if routeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", routeErr.Error())
		return
	}

	matches := route.Matches(nil, path)
	if matches {
		t.Errorf("the route was expected to NOT match")
		return
	}

	retHandler := route.Handler(nil, path)
	if retHandler != nil {
		t.Errorf("the returned handler was expected to be nil, handler returned")
		return
	}
}

func TestCreateRoute_withWriteRoute_userDoesNotHaveWriteAccess_Success(t *testing.T) {
	//variables:
	saveTrxFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (TransactionResponse, error) {
		return nil, nil
	}

	//execute:
	rols := roles.SDKFunc.Create()
	usrs := users.SDKFunc.Create()
	roleKey := "video-update-role-01"
	from := crypto.SDKFunc.GenPK().PublicKey()
	patternAsString := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	handler := createHandlerWithSaveTransactionFn(saveTrxFn)
	path := fmt.Sprintf("/videos/%s", "6adbfdfc-bb7d-4236-96d6-96d1688a2441")

	//execute:
	route, routeErr := createRoute(roleKey, rols, usrs, patternAsString, handler)
	if routeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", routeErr.Error())
		return
	}

	matches := route.Matches(from, path)
	if matches {
		t.Errorf("the route was expected to NOT match")
		return
	}

	retHandler := route.Handler(from, path)
	if retHandler != nil {
		t.Errorf("the returned handler was expected to be nil, handler returned")
		return
	}
}

func TestCreateRoute_withWriteRoute_userHasWriteAccess_Success(t *testing.T) {
	//variables:
	saveTrxFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (TransactionResponse, error) {
		return nil, nil
	}

	//execute:
	rols := roles.SDKFunc.Create()
	usrs := users.SDKFunc.Create()
	from := crypto.SDKFunc.GenPK().PublicKey()
	roleKey := "video-update-role-01"
	patternAsString := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	rolePatternAsString := "/videos/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})"
	handler := createHandlerWithSaveTransactionFn(saveTrxFn)
	path := fmt.Sprintf("/videos/%s", "6adbfdfc-bb7d-4236-96d6-96d1688a2441")

	// add the from user to the users list:
	usrs.Insert(from)

	//create the role with the from user:
	rols.Add(roleKey, from)

	// enable the role to write on the given url pattern:
	rols.EnableWriteAccess(roleKey, rolePatternAsString)

	//execute:
	route, routeErr := createRoute(roleKey, rols, usrs, patternAsString, handler)
	if routeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", routeErr.Error())
		return
	}

	matches := route.Matches(from, path)
	if !matches {
		t.Errorf("the route was expected to match")
		return
	}

	retHandler := route.Handler(from, path)
	if retHandler == nil {
		t.Errorf("the returned handler was expected to be valid, nil returned")
		return
	}
}

func TestCreateRouter_Success(t *testing.T) {

	//variables:
	rols := roles.SDKFunc.Create()
	usrs := users.SDKFunc.Create()
	roleKey := "video-update-role-01"

	// first route:
	firstQueryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		qr, qrErr := createEmptyQueryResponse(IsSuccessful, "first")
		return qr, qrErr
	}

	firstPatternAsString := "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	firstHandler := createHandlerWithQueryFn(firstQueryFn)
	firstRoute, firstRouteErr := createRoute(roleKey, rols, usrs, firstPatternAsString, firstHandler)
	if firstRouteErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstRouteErr.Error())
		return
	}

	// second route:
	secondQueryFn := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (QueryResponse, error) {
		qr, qrErr := createEmptyQueryResponse(IsSuccessful, "second")
		return qr, qrErr
	}

	secondPatternAsString := "/profiles/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>"
	secondHandler := createHandlerWithQueryFn(secondQueryFn)
	secondRoute, secondRouteErr := createRoute(roleKey, rols, usrs, secondPatternAsString, secondHandler)
	if secondRouteErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondRouteErr.Error())
		return
	}

	//execute:
	rtes := map[int][]Route{
		Retrieve: []Route{
			firstRoute,
			secondRoute,
		},
	}

	router := createRouter(rtes)
	firstPreparedHandler := router.Route(nil, "/videos/70de0f1a-0623-4bf6-ac6c-384f56321ec0", Retrieve)
	secondPrepatedHanlder := router.Route(nil, "/profiles/70de0f1a-0623-4bf6-ac6c-384f56321ec0", Retrieve)
	firstInvalidRouteHandler := router.Route(nil, "/this-is-invalid/70de0f1a-0623-4bf6-ac6c-384f56321ec0", Retrieve)
	secondInvalidRouteHandler := router.Route(nil, "/profiles/70de0f1a-0623-4bf6-ac6c-384f56321ec0", Save)

	if firstPreparedHandler == nil {
		t.Errorf("the returned prepared handler was expected to be valid, nil returned")
		return
	}

	if secondPrepatedHanlder == nil {
		t.Errorf("the returned prepared handler was expected to be valid, nil returned")
		return
	}

	if firstInvalidRouteHandler != nil {
		t.Errorf("the returned prepared handler was expected to be nil")
		return
	}

	if secondInvalidRouteHandler != nil {
		t.Errorf("the returned prepared handler was expected to be nil")
		return
	}

	firstFn := firstPreparedHandler.Handler().Query()
	firstResult, _ := firstFn(nil, nil, "", nil, nil)
	firstLog := firstResult.Log()
	if firstLog != "first" {
		t.Errorf("the returned log was invalid.  Expected: %s, Returned: %s", "first", firstLog)
	}

	secondFn := secondPrepatedHanlder.Handler().Query()
	secondResult, _ := secondFn(nil, nil, "", nil, nil)
	secondLog := secondResult.Log()
	if secondLog != "second" {
		t.Errorf("the returned log was invalid.  Expected: %s, Returned: %s", "first", secondLog)
	}
}
