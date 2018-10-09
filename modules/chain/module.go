package chain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	tendermint "github.com/xmnservices/xmnsuite/blockchains/tendermint"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	datastore_module "github.com/xmnservices/xmnsuite/modules/datastore"
	"github.com/xmnservices/xmnsuite/routers"
	lua "github.com/yuin/gopher-lua"
)

const luaChain = "chain"
const luaApplication = "app"
const luaRouter = "router"
const luaRoute = "route"

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

type module struct {
	context     *lua.LState
	dbPath      string
	port        int
	instanceID  *uuid.UUID
	rootPubKeys []crypto.PublicKey
	nodePK      tcrypto.PrivKey
	ds          datastore_module.Datastore
	ch          *chain
}

func createModule(
	context *lua.LState,
	dbPath string,
	port int,
	instanceID *uuid.UUID,
	rootPubKeys []crypto.PublicKey,
	nodePK tcrypto.PrivKey,
	ds datastore_module.Datastore,
) Chain {
	out := module{
		context:     context,
		dbPath:      dbPath,
		port:        port,
		instanceID:  instanceID,
		rootPubKeys: rootPubKeys,
		nodePK:      nodePK,
		ds:          ds,
		ch:          nil,
	}

	out.register()

	return &out
}

func (app *module) register() {
	// preload chain:
	app.context.PreloadModule("chain", func(context *lua.LState) int {
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
		}

		ntable := context.NewTable()
		context.SetFuncs(ntable, methods)
		context.Push(ntable)

		return 1
	})
}

func (app *module) registerChain(context *lua.LState) int {
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

func (app *module) registerApp(context *lua.LState) int {
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

func (app *module) registerRouter(context *lua.LState) int {
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

func (app *module) registerRoute(context *lua.LState) int {
	// create a new route instance:
	newRoute := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount != 3 {
			l.ArgError(1, "the exists func expected 3 parameters")
			return 1
		}

		routeTypeMapping := map[string]int{
			"retrieve": routers.Retrieve,
			"save":     routers.Save,
			"delete":   routers.Delete,
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

// Spawn spawns a new blockchain node.  A lua script with the blockchain contract must be executed first
func (app *module) Spawn() (applications.Node, error) {
	// make sure the chain is set:
	if app.ch == nil {
		return nil, errors.New("the chain has not been loaded")
	}

	// create the router data store:
	routerDS := datastore.SDKFunc.Create()

	appsSlice := []applications.Application{}
	for _, oneApp := range app.ch.apps {
		// create the route params:
		rteParams := []routers.CreateRouteParams{}
		for _, oneRte := range oneApp.router.rtes {
			var saveTrx routers.SaveTransactionFn
			if oneRte.saveTrx != nil {
				luaSaveTrxFn := oneRte.saveTrx
				saveTrx = func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinaryBare(from)
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
					sigAsString := sig.String()

					// call the func and return the value:
					return callLuaTrxFunc(
						luaSaveTrxFn,
						app.context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(dataAsString),
						lua.LString(sigAsString),
					)
				}
			}

			var delTrx routers.DeleteTransactionFn
			if oneRte.delTrx != nil {
				luaDelTrxFn := oneRte.delTrx
				delTrx = func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.TransactionResponse, error) {
					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinaryBare(from)
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
					sigAsString := sig.String()

					// call the func and return the value:
					return callLuaTrxFunc(
						luaDelTrxFn,
						app.context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(sigAsString),
					)
				}
			}

			var queryTrx routers.QueryFn
			if oneRte.queryTrx != nil {
				luaQueryFn := oneRte.queryTrx
				queryTrx = func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
					//replace the datastore:
					app.replaceDS(store)

					// from:
					fromAsBytes, fromAsBytesErr := cdc.MarshalBinaryBare(from)
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
					sigAsString := sig.String()

					// call the func and return the value:
					return callLuaQueryFunc(
						luaQueryFn,
						app.context,
						lua.LString(pubKeyAsString),
						lua.LString(path),
						&luaParams,
						lua.LString(sigAsString),
					)
				}
			}

			rteParams = append(rteParams, routers.CreateRouteParams{
				Pattern:  oneRte.pattern,
				SaveTrx:  saveTrx,
				DelTrx:   delTrx,
				QueryTrx: queryTrx,
			})
		}

		// setup the router role key:
		routerRoleKey := fmt.Sprintf("router-version-%s", oneApp.version)

		// add the root users on the routes:
		for _, onePubKey := range app.rootPubKeys {
			routerDS.Users().Insert(onePubKey)
			routerDS.Roles().Add(routerRoleKey, onePubKey)
			routerDS.Roles().EnableWriteAccess(routerRoleKey, ".*")
		}

		// create one application and put it in the list:
		appsSlice = append(appsSlice, applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
			Namespace:      app.ch.namespace,
			Name:           app.ch.name,
			ID:             app.instanceID,
			FromBlockIndex: int64(oneApp.beginIndex),
			ToBlockIndex:   int64(oneApp.endIndex),
			Version:        oneApp.version,
			DirPath:        app.dbPath,
			RouterParams: routers.CreateRouterParams{
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
		ID:        app.instanceID,
		PrivKey:   app.nodePK,
	})

	// create the blockchain service:
	blkChainService := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: app.dbPath,
	})

	// save the blockchain:
	saveBlkChainErr := blkChainService.Save(blkChain)
	if saveBlkChainErr != nil {
		return nil, saveBlkChainErr
	}

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(app.port, app.dbPath, blkChain, apps)
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

func callLuaQueryFunc(fn *lua.LFunction, context *lua.LState, args ...lua.LValue) (routers.QueryResponse, error) {
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

		return routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code:  code,
			Log:   log.String(),
			Key:   key.String(),
			Value: valueAsBytes,
		}), nil
	}

	return nil, errors.New("the query response is not a valid table")
}

func callLuaTrxFunc(fn *lua.LFunction, context *lua.LState, args ...lua.LValue) (routers.TransactionResponse, error) {
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

			return routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
				Code:    code,
				Log:     log.String(),
				GazUsed: int64(gazUsed),
				Tags:    tags,
			}), nil
		}

		return routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
			Code: code,
			Log:  log.String(),
		}), nil
	}

	return nil, errors.New("the transaction response is not a valid table")
}

func (app *module) replaceDS(store datastore.DataStore) *module {
	app.ds.Replace(store)
	return app
}
