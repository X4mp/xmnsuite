package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	applications "github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	tendermint "github.com/XMNBlockchain/datamint/tendermint"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type core struct {
	crypto *XCrypto
	keys   *XKeys
	tables *XTables
	users  *XUsers
	roles  *XRoles
	route  *XRoute
	router *XRouter
	app    *XApp
	chain  *XChain
}

func execute(dbPath string, instanceID *uuid.UUID, rootPubKeys []crypto.PubKey, nodePK crypto.PrivKey, store datastore.DataStore, context *lua.LState, scriptPath string) (applications.Node, error) {

	// create the core:
	core := createCore(context, store)

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		return nil, doFileErr
	}

	// make sure the chain is set:
	if core.chain.chain == nil {
		return nil, errors.New("the chain has not been loaded")
	}

	// create the router data store:
	routerDS := datastore.SDKFunc.Create()

	appsSlice := []applications.Application{}
	for _, oneApp := range core.chain.chain.apps {
		// create the route params:
		rteParams := []applications.CreateRouteParams{}
		for _, oneRte := range oneApp.routerParams.rtes {
			var saveTrx applications.SaveTransactionFn
			if oneRte.saveTrx != nil {
				luaSaveTrxFn := oneRte.saveTrx
				saveTrx = func(store datastore.DataStore, from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) (applications.TransactionResponse, error) {

					//replace the datastore:
					core.replaceDS(store)

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
					core.replaceDS(store)

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
					core.replaceDS(store)

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
			DataStore:      store,
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
		Namespace: core.chain.chain.namespace,
		Name:      core.chain.chain.name,
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

func (app *core) replaceDS(store datastore.DataStore) *core {
	app.tables.replaceObjects(store.Objects())
	return app
}

func createCore(l *lua.LState, store datastore.DataStore) *core {

	// preload JSON:
	luajson.Preload(l)

	// crypto:
	crypto := CreateXCrypto(l)

	// datastore:
	keys := CreateXKeys(l)
	tables := CreateXTables(l, store.Objects())

	// roles and users:
	users := CreateXUsers(l)
	roles := CreateXRoles(l)

	// router + application:
	route := CreateXRoute(l)
	router := CreateXRouter(l)
	app := CreateXApp(l)
	chain := CreateXChain(l)

	out := core{
		crypto: crypto,
		keys:   keys,
		tables: tables,
		users:  users,
		roles:  roles,
		route:  route,
		router: router,
		app:    app,
		chain:  chain,
	}

	return &out
}
