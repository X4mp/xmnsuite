package core

import (
	lua "github.com/yuin/gopher-lua"
)

const luaChain = "xchain"

type chainParams struct {
	namespace string
	name      string
	apps      []*appParams
}

// XChain represents the xchain instance
type XChain struct {
	l     *lua.LState
	chain *chainParams
}

// CreateXChain creates a new XChain instance:
func CreateXChain(l *lua.LState) *XChain {
	// create the instance:
	out := XChain{
		l:     l,
		chain: nil,
	}

	//registers the xcrypto module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XChain) register() {
	// convert the table argument to a chain:
	fromTableToChainFn := func(l *lua.LState) (*chainParams, error) {
		tb := l.ToTable(1)
		apps := []*appParams{}
		namespace := tb.RawGet(lua.LString("namespace"))
		name := tb.RawGet(lua.LString("name"))
		if rawApps, ok := tb.RawGet(lua.LString("apps")).(*lua.LTable); ok {
			rawApps.ForEach(func(key lua.LValue, rawApp lua.LValue) {
				if oneAppUD, ok := rawApp.(*lua.LUserData); ok {
					if oneApp, ok := oneAppUD.Value.(*appParams); ok {
						apps = append(apps, oneApp)
					}

				}
			})

		}

		return &chainParams{
			namespace: namespace.String(),
			name:      name.String(),
			apps:      apps,
		}, nil
	}

	// loadChain a loads apps into the chain:
	loadChain := func(l *lua.LState) int {

		if app.chain != nil {
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

		// add the chain params to the XChain:
		app.chain = chain

		// set the value:
		ud.Value = chain

		l.SetMetatable(ud, l.GetTypeMetatable(luaChain))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaChain)
	app.l.SetGlobal(luaChain, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadChain))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
