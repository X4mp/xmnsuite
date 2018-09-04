package core

import (
	"encoding/hex"
	"errors"
	"fmt"

	applications "github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	tendermint "github.com/XMNBlockchain/datamint/tendermint"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	lua "github.com/yuin/gopher-lua"
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

func execute(dbPath string, nodePK crypto.PrivKey, store datastore.DataStore, context *lua.LState, scriptPath string) (applications.Node, error) {

	// create the core:
	core := createCore(context)

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		return nil, doFileErr
	}

	// make sure the chain is set:
	if core.chain.chain == nil {
		return nil, errors.New("the chain has not been loaded")
	}

	appsSlice := []applications.Application{}
	for _, oneApp := range core.chain.chain.apps {
		// create the route params:
		rteParams := []applications.CreateRouteParams{}
		for _, oneRte := range oneApp.routerParams.rtes {

			var saveTrx applications.SaveTransactionFn
			if oneRte.saveTrx != nil {
				saveTrx = func(from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) (applications.TransactionResponse, error) {

					luaP := lua.P{
						Fn:      oneRte.saveTrx,
						NRet:    5,
						Protect: true,
					}

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

					// call the func:
					callErr := context.CallByParam(luaP, lua.LString(pubKeyAsString), lua.LString(path), &luaParams, lua.LString(dataAsString), lua.LString(sigAsString))
					if callErr != nil {
						return nil, callErr
					}

					// retrieve the returned value:
					luaRespTable := context.Get(-1)

					// remove the receives value:
					context.Pop(1)

					//create the transaction response:
					fmt.Printf("\n\n %v\n\n", luaRespTable)

					return nil, nil
				}
			}

			var delTrx applications.DeleteTransactionFn
			if oneRte.delTrx != nil {
				delTrx = func(from crypto.PubKey, path string, params map[string]string, sig []byte) (applications.TransactionResponse, error) {
					return nil, nil
				}
			}

			var queryTrx applications.QueryFn
			if oneRte.queryTrx != nil {
				queryTrx = func(from crypto.PubKey, path string, params map[string]string, sig []byte) (applications.QueryResponse, error) {
					return nil, nil
				}
			}

			rteParams = append(rteParams, applications.CreateRouteParams{
				Pattern:  oneRte.pattern,
				SaveTrx:  saveTrx,
				DelTrx:   delTrx,
				QueryTrx: queryTrx,
			})
		}

		// create one application and put it in the list:
		appsSlice = append(appsSlice, applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
			FromBlockIndex: int64(oneApp.beginIndex),
			ToBlockIndex:   int64(oneApp.endIndex),
			Version:        oneApp.version,
			DataStore:      store,
			RouterParams: applications.CreateRouterParams{
				DataStore:  store,
				RoleKey:    fmt.Sprintf("router-version-%s", oneApp.version),
				RtesParams: rteParams,
			},
		}))

	}

	// create the router data store:
	routerRoleKey := "router-role-key"
	routerDS := datastore.SDKFunc.Create()

	// variables - These should be in the lua script:
	namespace := "testapp"
	name := "MyTestApp"
	id := uuid.NewV4()
	fromPrivKey := ed25519.GenPrivKey()
	fromPubKey := fromPrivKey.PubKey()

	// add the users on the routes.  This should be in the lua script:
	routerDS.Users().Insert(fromPubKey)
	routerDS.Roles().Add(routerRoleKey, fromPubKey)
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/messages")

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: appsSlice,
	})

	// create the blockchain:
	blkChain := tendermint.SDKFunc.CreateBlockchain(tendermint.CreateBlockchainParams{
		Namespace: namespace,
		Name:      name,
		ID:        &id,
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
	defer node.Stop()

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		return nil, startNodeErr
	}

	return node, nil
}

func createCore(l *lua.LState) *core {
	// crypto:
	crypto := CreateXCrypto(l)

	// datastore:
	keys := CreateXKeys(l)
	tables := CreateXTables(l)

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
