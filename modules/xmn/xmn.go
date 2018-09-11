package xmn

import (
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	applications "github.com/xmnservices/xmnsuite/applications"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/roles"
	tendermint "github.com/xmnservices/xmnsuite/tendermint"
	"github.com/xmnservices/xmnsuite/users"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// core
const luaChain = "chain"
const luaApplication = "app"
const luaRouter = "router"
const luaRoute = "route"

// datastore
const luaTables = "tables"
const luaUsers = "users"
const luaRoles = "roles"
const luaKey = "keys"

// crypto:
const luaPrivKey = "privkey"

// sdk:
const luaResourcePointer = "rpointer"
const luaResource = "resource"
const luaTrxResponse = "trxresponse"
const luaQueryResponse = "queryresponse"

type chain struct {
	namespace string
	name      string
	apps      []*application
}

type application struct {
	version    string
	beginIndex int
	endIndex   int
	router     *router
}

type router struct {
	key  string
	rtes []*route
}

type route struct {
	pattern  string
	saveTrx  *lua.LFunction
	delTrx   *lua.LFunction
	queryTrx *lua.LFunction
}

type privKey struct {
	pk crypto.PrivKey
}

type resourcePointer struct {
	ptr applications.ResourcePointer
}

type resource struct {
	ptr applications.Resource
}

type xmn struct {
	ch     *chain
	ds     datastore.DataStore
	tables objects.Objects
	usrs   users.Users
	rols   roles.Roles
	k      keys.Keys
	cl     applications.Client
}

func createXMN(ds datastore.DataStore) *xmn {
	out := xmn{
		ch:     nil,
		ds:     ds,
		tables: ds.Objects(),
		usrs:   ds.Users(),
		rols:   ds.Roles(),
		k:      ds.Keys(),
		cl:     nil,
	}

	return &out
}

func createXMNWithClient(ds datastore.DataStore, client applications.Client) *xmn {
	out := xmn{
		ch:     nil,
		ds:     ds,
		tables: ds.Objects(),
		usrs:   ds.Users(),
		rols:   ds.Roles(),
		k:      ds.Keys(),
		cl:     client,
	}

	return &out
}

func (app *xmn) register(context *lua.LState) {
	// preload JSON:
	luajson.Preload(context)

	// preload XMN:
	context.PreloadModule("xmn", func(context *lua.LState) int {

		methods := map[string]lua.LGFunction{
			"chain": func(context *lua.LState) int {
				return app.registerChain(context)
			},
			"app": func(context *lua.LState) int {
				return app.registerApp(context)
			},
			"router": func(context *lua.LState) int {
				return app.registerRouter(context)
			},
			"route": func(context *lua.LState) int {
				return app.registerRoute(context)
			},
			"service": func(context *lua.LState) int {
				return app.registerService(context)
			},
		}

		ntable := context.NewTable()
		context.SetFuncs(ntable, methods)
		context.Push(ntable)

		// datastore:
		app.registerTables(context)
		app.registerTables(context)
		app.registerUsers(context)
		app.registerRoles(context)
		app.registerKeys(context)

		// crypto:
		app.registerPrivKey(context)

		// sdk:
		app.registerResourcePointer(context)
		app.registerResource(context)
		app.registerQueryResponse(context)
		app.registerClientTrxResponse(context)

		return 1
	})
}

func (app *xmn) registerChain(context *lua.LState) int {
	// convert the table argument to a chain:
	fromTableToChainFn := func(l *lua.LState) (*chain, error) {
		tb := l.ToTable(1)
		apps := []*application{}
		namespace := tb.RawGet(lua.LString("namespace"))
		name := tb.RawGet(lua.LString("name"))
		if rawApps, ok := tb.RawGet(lua.LString("apps")).(*lua.LTable); ok {
			rawApps.ForEach(func(key lua.LValue, rawApp lua.LValue) {
				if oneAppUD, ok := rawApp.(*lua.LUserData); ok {
					if oneApp, ok := oneAppUD.Value.(*application); ok {
						apps = append(apps, oneApp)
					}

				}
			})

		}

		return &chain{
			namespace: namespace.String(),
			name:      name.String(),
			apps:      apps,
		}, nil
	}

	// loadChain a loads apps into the chain:
	loadChain := func(l *lua.LState) int {

		if app.ch != nil {
			l.ArgError(1, "the chain has already been loaded")
			return 1
		}

		ud := l.NewUserData()

		amount := l.GetTop()
		if amount < 1 {
			l.ArgError(1, "the new function was expected to have at least 1 parameter")
			return 1
		}

		chain, chainErr := fromTableToChainFn(l)
		if chainErr != nil {
			l.ArgError(1, "the passed table argument is invalid")
			return 1
		}

		// add the chain params to the chain:
		app.ch = chain

		// set the value:
		ud.Value = chain

		l.SetMetatable(ud, l.GetTypeMetatable(luaChain))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"load": loadChain,
	}

	ntable := context.NewTable()
	context.SetFuncs(ntable, methods)
	context.Push(ntable)

	// return
	return 1
}

func (app *xmn) registerApp(context *lua.LState) int {
	// convert the table argument to an application:
	fromTableToAppFn := func(l *lua.LState) (*application, error) {
		tb := l.ToTable(1)
		version := tb.RawGet(lua.LString("version"))
		beginIndex := tb.RawGet(lua.LString("beginBlockIndex"))
		endIndex := tb.RawGet(lua.LString("endBlockIndex"))
		rterTable := tb.RawGet(lua.LString("router")).(*lua.LUserData)
		if router, ok := rterTable.Value.(*router); ok {

			beginIndexAsInt, beginIndexAsIntErr := strconv.Atoi(beginIndex.String())
			if beginIndexAsIntErr != nil {
				str := fmt.Sprintf("the given beginIndex (%d) is not a valid integer", beginIndex)
				return nil, errors.New(str)
			}

			endIndexAsInt, endIndexAsIntErr := strconv.Atoi(endIndex.String())
			if endIndexAsIntErr != nil {
				str := fmt.Sprintf("the given beginIndex (%d) is not a valid integer", beginIndex)
				return nil, errors.New(str)
			}

			return &application{
				version:    version.String(),
				beginIndex: beginIndexAsInt,
				endIndex:   endIndexAsInt,
				router:     router,
			}, nil
		}

		return nil, errors.New("the router param is invalid")
	}

	// create a new app instance:
	newApp := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount != 1 {
			l.ArgError(1, "the new function was expected to have 1 parameter")
			return 1
		}

		app, appErr := fromTableToAppFn(l)
		if appErr != nil {
			str := fmt.Sprintf("the passed table argument is invalid: %s", appErr.Error())
			l.ArgError(1, str)
			return 1
		}

		// set the value:
		ud.Value = app

		l.SetMetatable(ud, l.GetTypeMetatable(luaApplication))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"new": newApp,
	}

	ntable := context.NewTable()
	context.SetFuncs(ntable, methods)
	context.Push(ntable)

	// return:
	return 1
}

func (app *xmn) registerRouter(context *lua.LState) int {
	// convert the table argument to a router:
	fromTableToRouterFn := func(l *lua.LState) *router {
		tb := l.ToTable(1)

		routes := []*route{}
		key := tb.RawGet(lua.LString("key"))
		if rawRoutes, ok := tb.RawGet(lua.LString("routes")).(*lua.LTable); ok {
			rawRoutes.ForEach(func(key lua.LValue, rawRoute lua.LValue) {
				if oneRouteUD, ok := rawRoute.(*lua.LUserData); ok {
					if oneRoute, ok := oneRouteUD.Value.(*route); ok {
						routes = append(routes, oneRoute)
					}

				}
			})

		}

		return &router{
			key:  key.String(),
			rtes: routes,
		}
	}

	// create a new router instance:
	newRouter := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount != 1 {
			l.ArgError(1, "the new function was expected to have 1 parameter")
			return 1
		}

		router := fromTableToRouterFn(l)
		if router == nil {
			l.ArgError(1, "the passed table argument is invalid")
			return 1
		}

		// set the value:
		ud := l.NewUserData()
		ud.Value = router

		l.SetMetatable(ud, l.GetTypeMetatable(luaRouter))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"new": newRouter,
	}

	ntable := context.NewTable()
	context.SetFuncs(ntable, methods)
	context.Push(ntable)

	// return:
	return 1
}

func (app *xmn) registerRoute(context *lua.LState) int {
	// create a new route instance:
	newRoute := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount != 3 {
			l.ArgError(1, "the exists func expected 3 parameters")
			return 1
		}

		routeTypeMapping := map[string]int{
			"retrieve": applications.Retrieve,
			"save":     applications.Save,
			"delete":   applications.Delete,
		}

		rteTypeAsString := l.CheckString(1)
		if _, ok := routeTypeMapping[rteTypeAsString]; !ok {
			l.ArgError(2, "the passed route type is invalid")
			return 1
		}

		patternAsString := l.CheckString(2)
		luaHandlr := l.CheckFunction(3)

		newRte := route{
			pattern:  patternAsString,
			saveTrx:  nil,
			delTrx:   nil,
			queryTrx: nil,
		}

		// if the handler is retrieve:
		if rteTypeAsString == "retrieve" {
			if luaHandlr.Proto.NumParameters != 4 {
				l.ArgError(1, "the retrieve func handler is expected to have 4 parameters")
				return 3
			}

			newRte.queryTrx = luaHandlr
		}

		// if the handler is delete:
		if rteTypeAsString == "delete" {
			if luaHandlr.Proto.NumParameters != 4 {
				l.ArgError(1, "the delete func handler is expected to have 4 parameters")
				return 3
			}

			newRte.delTrx = luaHandlr
		}

		// if the handler is save:
		if rteTypeAsString == "save" {
			if luaHandlr.Proto.NumParameters != 5 {
				l.ArgError(1, "the save func handler is expected to have 5 parameters")
				return 3
			}

			newRte.saveTrx = luaHandlr
		}

		// set the value:
		ud.Value = &newRte

		l.SetMetatable(ud, l.GetTypeMetatable(luaRoute))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"new": newRoute,
	}

	ntable := context.NewTable()
	context.SetFuncs(ntable, methods)
	context.Push(ntable)

	// return:
	return 1
}

func (app *xmn) registerService(context *lua.LState) int {
	// transactFn executes a transaction on the service
	transactFn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount != 1 {
			l.ArgError(1, "the transact func expected 1 parameter")
			return 1
		}

		if app.cl == nil {
			l.RaiseError("the client must be set in the engine in order to use the service")
			return 1
		}

		tb := l.ToTable(1)
		luaSig := tb.RawGetString("sig")
		sig, sigErr := hex.DecodeString(luaSig.String())
		if sigErr != nil {
			l.ArgError(1, "the sig could not be decoded")
			return 1
		}

		params := applications.CreateTransactionRequestParams{
			Sig: sig,
		}

		luaRes := tb.RawGetString("resource")
		if luaRes.Type().String() != lua.LTNil.String() {
			if restUD, ok := luaRes.(*lua.LUserData); ok {
				if res, ok := restUD.Value.(*resource); ok {
					params.Res = res.ptr
				}
			}
		}

		luaResPtr := tb.RawGetString("rpointer")
		if luaResPtr.Type().String() != lua.LTNil.String() {
			if restPtrUD, ok := luaResPtr.(*lua.LUserData); ok {
				if resPtr, ok := restPtrUD.Value.(*resourcePointer); ok {
					params.Ptr = resPtr.ptr
				}
			}
		}

		// create the request:
		req := applications.SDKFunc.CreateTransactionRequest(params)

		// execte the request:
		resp, respErr := app.cl.Transact(req)
		if respErr != nil {
			str := fmt.Sprintf("there was an error while execting the transaction request: %s", respErr.Error())
			l.ArgError(1, str)
			return 1
		}

		// set the response:
		ud := l.NewUserData()
		ud.Value = resp

		l.SetMetatable(ud, l.GetTypeMetatable(luaTrxResponse))
		l.Push(ud)
		return 1
	}

	// queryFn executes a query on the service
	queryFn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount != 1 {
			l.ArgError(1, "the transact func expected 1 parameter")
			return 1
		}

		if app.cl == nil {
			l.RaiseError("the client must be set in the engine in order to use the service")
			return 1
		}

		tb := l.ToTable(1)
		luaSig := tb.RawGetString("sig")
		sig, sigErr := hex.DecodeString(luaSig.String())
		if sigErr != nil {
			l.ArgError(1, "the sig could not be decoded")
			return 1
		}

		params := applications.CreateQueryRequestParams{
			Sig: sig,
		}

		luaResPtr := tb.RawGetString("rpointer")
		if luaResPtr.Type().String() != lua.LTNil.String() {
			if restPtrUD, ok := luaResPtr.(*lua.LUserData); ok {
				if resPtr, ok := restPtrUD.Value.(*resourcePointer); ok {
					params.Ptr = resPtr.ptr
				}
			}
		}

		// create the request:
		req := applications.SDKFunc.CreateQueryRequest(params)

		// execte the request:
		resp, respErr := app.cl.Query(req)
		if respErr != nil {
			str := fmt.Sprintf("there was an error while execting the query request: %s", respErr.Error())
			l.ArgError(1, str)
			return 1
		}

		// set the response:
		ud := l.NewUserData()
		ud.Value = resp

		l.SetMetatable(ud, l.GetTypeMetatable(luaQueryResponse))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"transact": transactFn,
		"query":    queryFn,
	}

	ntable := context.NewTable()
	context.SetFuncs(ntable, methods)
	context.Push(ntable)

	// return:
	return 1
}

func (app *xmn) registerQueryResponse(context *lua.LState) int {
	checkFn := func(l *lua.LState) applications.QueryResponse {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(applications.QueryResponse); ok {
			return v
		}

		l.ArgError(1, "QueryResponse expected")
		return nil
	}

	// codeFn returns the query response code
	codeFn := func(l *lua.LState) int {
		resp := checkFn(l)
		l.Push(lua.LNumber(resp.Code()))
		return 1
	}

	// logFn returns the query response log
	logFn := func(l *lua.LState) int {
		resp := checkFn(l)
		l.Push(lua.LString(resp.Log()))
		return 1
	}

	// keyFn returns the query response key
	keyFn := func(l *lua.LState) int {
		resp := checkFn(l)
		l.Push(lua.LString(resp.Key()))
		return 1
	}

	// ValueFn returns the query response value
	ValueFn := func(l *lua.LState) int {
		resp := checkFn(l)
		l.Push(lua.LString(string(resp.Value())))
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"code":  codeFn,
		"log":   logFn,
		"key":   keyFn,
		"value": ValueFn,
	}

	mt := context.NewTypeMetatable(luaQueryResponse)

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))

	// return:
	return 1
}

func (app *xmn) registerClientTrxResponse(context *lua.LState) int {
	checkFn := func(l *lua.LState) applications.ClientTransactionResponse {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(applications.ClientTransactionResponse); ok {
			return v
		}

		l.ArgError(1, "ClientTransactionResponse expected")
		return nil
	}

	// codeFn returns the query response code
	codeFn := func(l *lua.LState) int {
		resp := checkFn(l)

		chk := resp.Check()
		if chk.Code() != applications.IsSuccessful {
			l.Push(lua.LNumber(chk.Code()))
			return 1
		}

		l.Push(lua.LNumber(resp.Transaction().Code()))
		return 1
	}

	// logFn returns the query response log
	logFn := func(l *lua.LState) int {
		resp := checkFn(l)

		chk := resp.Check()
		if chk.Code() != applications.IsSuccessful {
			l.Push(lua.LString(chk.Log()))
			return 1
		}

		l.Push(lua.LString(resp.Transaction().Log()))
		return 1
	}

	// gazUsedFn returns the amont of gaz used
	gazUsedFn := func(l *lua.LState) int {
		resp := checkFn(l)
		l.Push(lua.LNumber(resp.Transaction().GazUsed()))
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"code":    codeFn,
		"log":     logFn,
		"gazUsed": gazUsedFn,
	}

	mt := context.NewTypeMetatable(luaTrxResponse)

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))

	// return:
	return 1
}

func (app *xmn) registerTables(context *lua.LState) int {
	//gob register:
	gob.Register(map[string]interface{}{})

	//verifies that the given type is an object instance:
	checkFn := func(l *lua.LState) objects.Objects {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(objects.Objects); ok {
			return v
		}

		l.ArgError(1, "tables expected")
		return nil
	}

	// load the Objects instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.tables
		l.SetMetatable(ud, l.GetTypeMetatable(luaTables))
		l.Push(ud)
		return 1
	}

	//execute the save command on the objects instance:
	saveFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the save func expected at least 2 parameters")
			return 1
		}

		params := []*objects.ObjInKey{}
		for i := 2; i <= amount; i++ {
			oneObjInKey := objects.ObjInKey{}
			oneParam := l.CheckTable(i)
			oneParam.ForEach(func(name lua.LValue, value lua.LValue) {
				valueType := value.Type()
				nameAsString := name.String()
				if nameAsString == "table" && valueType == lua.LTTable {
					oneObjInKey.Obj = convertLTableToHashMap(value.(*lua.LTable))
				}

				if nameAsString == "key" && valueType == lua.LTString {
					oneObjInKey.Key = value.String()
				}
			})

			params = append(params, &oneObjInKey)
		}

		amountSaved := p.Save(params...)
		l.Push(lua.LNumber(amountSaved))
		return 1
	}

	//execute the retrieve command on the objects instance:
	retrieveFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		objInKey := objects.ObjInKey{
			Key: l.CheckString(2),
			Obj: new(map[string]interface{}),
		}

		p.Retrieve(&objInKey)
		mapResult := objInKey.Obj.(*map[string]interface{})
		ltable := convertHashMapToLTable(*mapResult)
		if ltable == nil {
			l.Push(lua.LNil)
			return 1
		}

		l.Push(ltable)
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p.Keys())(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkFn(l)
			return existsFn(p.Keys())(l)
		},
		"retrieve": retrieveFn,
		"search": func(l *lua.LState) int {
			p := checkFn(l)
			return searchFn(p.Keys())(l)
		},
		"save": saveFn,
		"delete": func(l *lua.LState) int {
			p := checkFn(l)
			return delFn(p.Keys())(l)
		},
	}

	mt := context.NewTypeMetatable(luaTables)
	context.SetGlobal(luaTables, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))

	// return:
	return 1
}

func (app *xmn) registerUsers(context *lua.LState) {
	//verifies that the given type is a Users instance:
	checkFn := func(l *lua.LState) users.Users {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(users.Users); ok {
			return v
		}

		l.ArgError(1, "users expected")
		return nil
	}

	// load the Users instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.usrs
		l.SetMetatable(ud, l.GetTypeMetatable(luaUsers))
		l.Push(ud)
		return 1
	}

	//execute the key command on the objects instance:
	keyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		key := p.Key(pubKey)
		l.Push(lua.LString(key))
		return 1
	}

	//execute the exists command on the objects instance:
	existsFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		exists := p.Exists(pubKey)
		l.Push(lua.LBool(exists))
		return 1
	}

	//execute the insert command on the objects instance:
	insertFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isInserted := p.Insert(pubKey)
		l.Push(lua.LBool(isInserted))
		return 1
	}

	//execute the delete command on the objects instance:
	deleteFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isDeleted := p.Delete(pubKey)
		l.Push(lua.LBool(isDeleted))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p.Objects().Keys())(l)
		},
		"key":    keyFn,
		"exists": existsFn,
		"insert": insertFn,
		"delete": deleteFn,
	}

	mt := context.NewTypeMetatable(luaUsers)
	context.SetGlobal(luaUsers, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) registerRoles(context *lua.LState) {
	//verifies that the given type is a Roles instance:
	checkFn := func(l *lua.LState) roles.Roles {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(roles.Roles); ok {
			return v
		}

		l.ArgError(1, "roles expected, received %d")
		return nil
	}

	// load the Roles instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.rols
		l.SetMetatable(ud, l.GetTypeMetatable(luaRoles))
		l.Push(ud)
		return 1
	}

	//execute the add command on the roles instance:
	addFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Add(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the del command on the roles instance:
	delFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Del(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the enableWriteAccess command on the roles instance:
	enableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.EnableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	//execute the disableWriteAccess command on the roles instance:
	disableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.DisableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	//execute the hasWriteAccess command on the roles instance:
	hasWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		keys := []string{}
		for i := 3; i <= amount; i++ {
			keys = append(keys, l.CheckString(i))
		}

		returnedKeys := p.HasWriteAccess(key, keys...)
		table := lua.LTable{}
		for _, oneKey := range returnedKeys {
			table.Append(lua.LString(oneKey))
		}

		l.Push(&table)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"add":                addFn,
		"del":                delFn,
		"enableWriteAccess":  enableWriteAccessFn,
		"disableWriteAccess": disableWriteAccessFn,
		"hasWriteAccess":     hasWriteAccessFn,
	}

	mt := context.NewTypeMetatable(luaRoles)
	context.SetGlobal(luaRoles, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) registerKeys(context *lua.LState) {
	//verifies that the given type is a keys instance:
	checkFn := func(l *lua.LState) keys.Keys {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(keys.Keys); ok {
			return v
		}

		l.ArgError(1, "keys expected")
		return nil
	}

	// load the Keys instance:
	loadKeys := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.k
		l.SetMetatable(ud, l.GetTypeMetatable(luaKey))
		l.Push(ud)
		return 1
	}

	// execute the retrieve command on the keys instance:
	retrieveFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount != 2 {
			l.ArgError(1, "the retrieve func expected 1 parameter")
			return 1
		}

		key := l.CheckString(2)
		value := p.Retrieve(key)
		if value == nil {
			l.Push(&lua.LNilType{})
			return 1
		}

		l.Push(lua.LString(value.(string)))
		return 1
	}

	//execute the save command on the keys instance:
	saveFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 3 {
			key := l.CheckString(2)
			value := l.CheckString(3)
			p.Save(key, value)
			return 0
		}

		l.ArgError(1, "the save func expected 2 parameters")
		return 1
	}

	// the keys methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p)(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkFn(l)
			return existsFn(p)(l)
		},
		"retrieve": retrieveFn,
		"search": func(l *lua.LState) int {
			p := checkFn(l)
			return searchFn(p)(l)
		},
		"save": saveFn,
		"delete": func(l *lua.LState) int {
			p := checkFn(l)
			return delFn(p)(l)
		},
	}

	mt := context.NewTypeMetatable(luaKey)
	context.SetGlobal(luaKey, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadKeys))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) registerResource(context *lua.LState) {

	//verifies that the given type is a Resource instance:
	checkFn := func(l *lua.LState) *resource {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*resource); ok {
			return v
		}

		l.ArgError(1, "resource expected")
		return nil
	}

	// create a new resource instance:
	newResource := func(l *lua.LState) int {
		ud := l.NewUserData()
		if l.GetTop() == 1 {
			dataTable := l.CheckTable(1)
			ptr := dataTable.RawGetString("pointer")
			data := dataTable.RawGetString("data")

			if ptrUD, ok := ptr.(*lua.LUserData); ok {
				if pointer, ok := ptrUD.Value.(*resourcePointer); ok {

					res := applications.SDKFunc.CreateResource(applications.CreateResourceParams{
						ResPtr: pointer.ptr,
						Data:   []byte(data.String()),
					})

					ud.Value = &resource{
						ptr: res,
					}

					l.SetMetatable(ud, l.GetTypeMetatable(luaResource))
					l.Push(ud)
					return 1
				}
			}

			l.ArgError(1, "the new func expected its parameter to be a table that contains a pointer keyname, that reference a resource pointer instance")
			return 1

		}

		l.ArgError(1, "the new func expected 1 parameter")
		return 1
	}

	// hashFn executes the hash fnc on the resource instance
	hashFn := func(l *lua.LState) int {
		res := checkFn(l)
		if res == nil {
			return 1
		}

		hashAsString := hex.EncodeToString(res.ptr.Hash())
		l.Push(lua.LString(hashAsString))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"hash": hashFn,
	}

	mt := context.NewTypeMetatable(luaResource)
	context.SetGlobal(luaResource, mt)

	// static attributes
	context.SetField(mt, "new", context.NewFunction(newResource))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) registerResourcePointer(context *lua.LState) {

	//verifies that the given type is a ResourcePointer instance:
	checkFn := func(l *lua.LState) *resourcePointer {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*resourcePointer); ok {
			return v
		}

		l.ArgError(1, "resource pointer expected")
		return nil
	}

	// create a new resource pointer instance:
	newResourcePointer := func(l *lua.LState) int {
		ud := l.NewUserData()
		if l.GetTop() == 1 {
			dataTable := l.CheckTable(1)
			from := dataTable.RawGetString("from")
			path := dataTable.RawGetString("path")

			// decode the pubkey:
			pubKeyAsBytes, pubKeyAsBytesErr := hex.DecodeString(from.String())
			if pubKeyAsBytesErr != nil {
				l.ArgError(1, "the from pubKey cold not be decoded is invalid")
				return 1
			}

			// unmarshal the bytes:
			newPubKey := new(ed25519.PubKeyEd25519)
			pubKerr := cdc.UnmarshalBinary(pubKeyAsBytes, newPubKey)
			if pubKerr != nil {
				l.ArgError(1, "the from pubKey bytes cannot be converted to a pubKey instance")
				return 1
			}

			// create the resource pointer:
			ptr := applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
				From: newPubKey,
				Path: path.String(),
			})

			ud.Value = &resourcePointer{
				ptr: ptr,
			}

			l.SetMetatable(ud, l.GetTypeMetatable(luaResourcePointer))
			l.Push(ud)
			return 1
		}

		l.ArgError(1, "the new func expected 1 parameter")
		return 1
	}

	// hashFn executes the hash fnc on the resource pointer instance
	hashFn := func(l *lua.LState) int {
		presource := checkFn(l)
		if presource == nil {
			return 1
		}

		hashAsString := hex.EncodeToString(presource.ptr.Hash())
		l.Push(lua.LString(hashAsString))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"hash": hashFn,
	}

	mt := context.NewTypeMetatable(luaResourcePointer)
	context.SetGlobal(luaResourcePointer, mt)

	// static attributes
	context.SetField(mt, "new", context.NewFunction(newResourcePointer))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) registerPrivKey(context *lua.LState) {
	//verifies that the given type is a Crypto instance:
	checkFn := func(l *lua.LState) *privKey {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*privKey); ok {
			return v
		}

		l.ArgError(1, "private key expected")
		return nil
	}

	// create a new crypto instance:
	newPrivKey := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = &privKey{
			pk: ed25519.GenPrivKey(),
		}

		l.SetMetatable(ud, l.GetTypeMetatable(luaPrivKey))
		l.Push(ud)
		return 1
	}

	//execute the pubKey command on the privkey instance:
	pubKeyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if p == nil {
			return 1
		}

		pubKeyAsBytes, pubKeyAsBytesErr := cdc.MarshalBinary(p.pk.PubKey())
		if pubKeyAsBytesErr != nil {
			l.ArgError(1, "the public key of the private key is invalid")
			return 1
		}

		pubKey := hex.EncodeToString(pubKeyAsBytes)
		l.Push(lua.LString(pubKey))
		return 1
	}

	//execute the signFn command on the privkey instance:
	signFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 2 {
			encodedHash := l.CheckString(2)
			msg, msgErr := hex.DecodeString(encodedHash)
			if msgErr != nil {
				str := fmt.Sprintf("the hash could not be decoded: %s", msgErr.Error())
				l.ArgError(1, str)
				return 1
			}

			sig, sigErr := p.pk.Sign(msg)
			if sigErr != nil {
				str := fmt.Sprintf("there was an error while signing the given message: %s", sigErr.Error())
				l.ArgError(1, str)
				return 1
			}

			encocedSig := hex.EncodeToString(sig)
			l.Push(lua.LString(encocedSig))
			return 1
		}

		l.ArgError(1, "the save func expected 1 parameter")
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"pubKey": pubKeyFn,
		"sign":   signFn,
	}

	mt := context.NewTypeMetatable(luaPrivKey)
	context.SetGlobal(luaPrivKey, mt)

	// static attributes
	context.SetField(mt, "new", context.NewFunction(newPrivKey))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *xmn) execute(
	context *lua.LState,
	dbPath string,
	instanceID *uuid.UUID,
	rootPubKeys []crypto.PubKey,
	nodePK crypto.PrivKey,
	scriptPath string,
) (applications.Node, error) {

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		return nil, doFileErr
	}

	// make sure the chain is set:
	if app.ch == nil {
		return nil, errors.New("the chain has not been loaded")
	}

	// create the router data store:
	routerDS := datastore.SDKFunc.Create()

	appsSlice := []applications.Application{}
	for _, oneApp := range app.ch.apps {
		// create the route params:
		rteParams := []applications.CreateRouteParams{}
		for _, oneRte := range oneApp.router.rtes {
			var saveTrx applications.SaveTransactionFn
			if oneRte.saveTrx != nil {
				luaSaveTrxFn := oneRte.saveTrx
				saveTrx = func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) (applications.TransactionResponse, error) {

					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinary(from)
					if fromAsBytesErr != nil {
						return nil, fromAsBytesErr
					}

					pubKeyAsString := hex.EncodeToString(fromAsBytes)

					// params:
					luaParams := lua.LTable{}
					for keyname, value := range params {
						luaParams.RawSet(lua.LString(keyname), lua.LString(value))
					}

					// json data as string:
					dataAsString := string(data)

					// sig:
					sigAsString := hex.EncodeToString(sig)

					// call the func and return the value:
					return callLuaTrxFunc(
						luaSaveTrxFn,
						context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(dataAsString),
						lua.LString(sigAsString),
					)
				}
			}

			var delTrx applications.DeleteTransactionFn
			if oneRte.delTrx != nil {
				luaDelTrxFn := oneRte.delTrx
				delTrx = func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) (applications.TransactionResponse, error) {
					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinary(from)
					if fromAsBytesErr != nil {
						return nil, fromAsBytesErr
					}

					pubKeyAsString := hex.EncodeToString(fromAsBytes)

					// params:
					luaParams := lua.LTable{}
					for keyname, value := range params {
						luaParams.RawSet(lua.LString(keyname), lua.LString(value))
					}

					// sig:
					sigAsString := hex.EncodeToString(sig)

					// call the func and return the value:
					return callLuaTrxFunc(
						luaDelTrxFn,
						context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(sigAsString),
					)
				}
			}

			var queryTrx applications.QueryFn
			if oneRte.queryTrx != nil {
				luaQueryFn := oneRte.queryTrx
				queryTrx = func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, sig []byte) (applications.QueryResponse, error) {
					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinary(from)
					if fromAsBytesErr != nil {
						return nil, fromAsBytesErr
					}

					pubKeyAsString := hex.EncodeToString(fromAsBytes)

					// params:
					luaParams := lua.LTable{}
					for keyname, value := range params {
						luaParams.RawSet(lua.LString(keyname), lua.LString(value))
					}

					// sig:
					sigAsString := hex.EncodeToString(sig)

					// call the func and return the value:
					return callLuaQueryFunc(
						luaQueryFn,
						context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(sigAsString),
					)
				}
			}

			rteParams = append(rteParams, applications.CreateRouteParams{
				Pattern:  oneRte.pattern,
				SaveTrx:  saveTrx,
				DelTrx:   delTrx,
				QueryTrx: queryTrx,
			})
		}

		// setup the router role key:
		routerRoleKey := fmt.Sprintf("router-version-%s", oneApp.version)

		// add the root users on the routes:
		for _, onePubKey := range rootPubKeys {
			routerDS.Users().Insert(onePubKey)
			routerDS.Roles().Add(routerRoleKey, onePubKey)
			routerDS.Roles().EnableWriteAccess(routerRoleKey, "/messages")
			routerDS.Roles().EnableWriteAccess(routerRoleKey, "/messages/[a-z0-9-]+")
		}

		// create one application and put it in the list:
		appsSlice = append(appsSlice, applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
			FromBlockIndex: int64(oneApp.beginIndex),
			ToBlockIndex:   int64(oneApp.endIndex),
			Version:        oneApp.version,
			DataStore:      app.ds,
			RouterParams: applications.CreateRouterParams{
				DataStore:  routerDS,
				RoleKey:    routerRoleKey,
				RtesParams: rteParams,
			},
		}))

	}

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: appsSlice,
	})

	// create the blockchain:
	blkChain := tendermint.SDKFunc.CreateBlockchain(tendermint.CreateBlockchainParams{
		Namespace: app.ch.namespace,
		Name:      app.ch.name,
		ID:        instanceID,
		PrivKey:   nodePK,
	})

	// create the blockchain service:
	blkChainService := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: dbPath,
	})

	// save the blockchain:
	saveBlkChainErr := blkChainService.Save(blkChain)
	if saveBlkChainErr != nil {
		return nil, saveBlkChainErr
	}

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService(tendermint.CreateApplicationServiceParams{
		RootDir:  dbPath,
		BlkChain: blkChain,
		Apps:     apps,
	})

	// spawn the node:
	node, nodeErr := appService.Spawn()
	if nodeErr != nil {
		return nil, nodeErr
	}

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		return nil, startNodeErr
	}

	return node, nil
}

func callLuaQueryFunc(fn *lua.LFunction, context *lua.LState, args ...lua.LValue) (applications.QueryResponse, error) {
	luaP := lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}

	// call the func:
	callErr := context.CallByParam(luaP, args...)
	if callErr != nil {
		return nil, callErr
	}

	// retrieve the returned value:
	value := context.Get(-1)
	context.Pop(1)
	if luaRespTable, ok := value.(*lua.LTable); ok {
		// fetch the data:
		codeAsLua := luaRespTable.RawGetString("code")
		log := luaRespTable.RawGetString("log")
		key := luaRespTable.RawGetString("key")
		value := luaRespTable.RawGetString("value")

		code, codeErr := strconv.Atoi(codeAsLua.String())
		if codeErr != nil {
			str := fmt.Sprintf("the code (%s) in the return table is not a valid integer", codeAsLua.String())
			return nil, errors.New(str)
		}

		valueAsBytes := []byte(value.String())
		if value.Type() == lua.LNil.Type() {
			valueAsBytes = nil
		}

		return applications.SDKFunc.CreateQueryResponse(applications.CreateQueryResponseParams{
			Code:  code,
			Log:   log.String(),
			Key:   key.String(),
			Value: valueAsBytes,
		}), nil
	}

	return nil, errors.New("the query response is not a valid table")
}

func callLuaTrxFunc(fn *lua.LFunction, context *lua.LState, args ...lua.LValue) (applications.TransactionResponse, error) {
	luaP := lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}

	// call the func:
	callErr := context.CallByParam(luaP, args...)
	if callErr != nil {
		return nil, callErr
	}

	// retrieve the returned value:
	value := context.Get(-1)
	context.Pop(1)
	if luaRespTable, ok := value.(*lua.LTable); ok {
		// fetch the data:
		codeAsLua := luaRespTable.RawGetString("code")
		log := luaRespTable.RawGetString("log")
		gazUsedAsLua := luaRespTable.RawGetString("gazUsed")
		luaTags := luaRespTable.RawGetString("tags")

		code, codeErr := strconv.Atoi(codeAsLua.String())
		if codeErr != nil {
			str := fmt.Sprintf("the code (%s) in the return table is not a valid integer", codeAsLua.String())
			return nil, errors.New(str)
		}

		if gazUsedAsLua != lua.LNil && luaTags != lua.LNil {
			tags := map[string][]byte{}
			if rawTags, ok := luaTags.(*lua.LTable); ok {
				rawTags.ForEach(func(key lua.LValue, luaKeyValueTable lua.LValue) {
					if rawKeyValueTable, ok := luaKeyValueTable.(*lua.LTable); ok {
						tagKey := rawKeyValueTable.RawGetString("key")
						tagValue := rawKeyValueTable.RawGetString("value")
						tags[tagKey.String()] = []byte(tagValue.String())
					}
				})

			}

			gazUsed, gazUsedErr := strconv.Atoi(gazUsedAsLua.String())
			if gazUsedErr != nil {
				str := fmt.Sprintf("the gazUsed (%s) in the return table is not a valid integer", gazUsedAsLua.String())
				return nil, errors.New(str)
			}

			return applications.SDKFunc.CreateTransactionResponse(applications.CreateTransactionResponseParams{
				Code:    code,
				Log:     log.String(),
				GazUsed: int64(gazUsed),
				Tags:    tags,
			}), nil
		}

		return applications.SDKFunc.CreateTransactionResponse(applications.CreateTransactionResponseParams{
			Code: code,
			Log:  log.String(),
		}), nil
	}

	return nil, errors.New("the transaction response is not a valid table")
}

func (app *xmn) replaceDS(store datastore.DataStore) *xmn {
	app.tables = store.Objects()
	return app
}
