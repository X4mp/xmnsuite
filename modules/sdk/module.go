package sdk

import (
	"fmt"

	applications "github.com/xmnservices/xmnsuite/applications"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	lua "github.com/yuin/gopher-lua"
)

const luaResourcePointer = "rpointer"
const luaResource = "resource"
const luaTrxResponse = "trxresponse"
const luaQueryResponse = "queryresponse"

type resourcePointer struct {
	ptr applications.ResourcePointer
}

type resource struct {
	ptr applications.Resource
}

type module struct {
	context *lua.LState
	client  applications.Client
}

func createModule(context *lua.LState, client applications.Client) *module {
	out := module{
		context: context,
		client:  client,
	}

	out.register()

	return &out
}

func (app *module) register() {
	app.context.PreloadModule("sdk", func(context *lua.LState) int {

		methods := map[string]lua.LGFunction{
			"service": func(context *lua.LState) int {
				return app.registerService(context)
			},
		}

		ntable := context.NewTable()
		context.SetFuncs(ntable, methods)
		context.Push(ntable)

		app.registerResourcePointer(context)
		app.registerResource(context)
		app.registerQueryResponse(context)
		app.registerClientTrxResponse(context)

		return 1
	})
}

func (app *module) registerResourcePointer(context *lua.LState) {

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

			// unmarshal the bytes:
			newPubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
				PubKeyAsString: from.String(),
			})

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

		l.Push(lua.LString(presource.ptr.Hash()))
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

func (app *module) registerResource(context *lua.LState) {

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

		l.Push(lua.LString(res.ptr.Hash()))
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

func (app *module) registerService(context *lua.LState) int {
	// transactFn executes a transaction on the service
	transactFn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount != 1 {
			l.ArgError(1, "the transact func expected 1 parameter")
			return 1
		}

		tb := l.ToTable(1)
		params := applications.CreateTransactionRequestParams{
			Sig: crypto.SDKFunc.CreateSig(crypto.CreateSigParams{
				SigAsString: tb.RawGetString("sig").String(),
			}),
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
		resp, respErr := app.client.Transact(req)
		if respErr != nil {
			str := fmt.Sprintf("there was an error while executing the transaction request: %s", respErr.Error())
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

		if app.client == nil {
			l.RaiseError("the client must be set in the engine in order to use the service")
			return 1
		}

		tb := l.ToTable(1)
		params := applications.CreateQueryRequestParams{
			Sig: crypto.SDKFunc.CreateSig(crypto.CreateSigParams{
				SigAsString: tb.RawGetString("sig").String(),
			}),
		}

		luaResPtr := tb.RawGetString("rpointer")
		if luaResPtr.Type().String() != lua.LTNil.String() {
			if restPtrUD, ok := luaResPtr.(*lua.LUserData); ok {
				if resPtr, ok := restPtrUD.Value.(*resourcePointer); ok {
					params.Ptr = resPtr.ptr
				}
			}
		}

		if params.Ptr == nil {
			l.ArgError(1, "the params expected an rpointer value")
			return 1
		}

		// create the request:
		req := applications.SDKFunc.CreateQueryRequest(params)

		// execte the request:
		resp, respErr := app.client.Query(req)
		if respErr != nil {
			l.RaiseError("there was an error while executing the query request: %s", respErr.Error())
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

func (app *module) registerQueryResponse(context *lua.LState) int {
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

func (app *module) registerClientTrxResponse(context *lua.LState) int {
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
