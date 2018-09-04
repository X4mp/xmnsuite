package core

import (
	lua "github.com/yuin/gopher-lua"
)

const luaChain = "xchain"

type chainParams struct {
	apps []*appParams
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
	//verifies that the given type is a App instance:
	checkAppFn := func(l *lua.LState, index int) *appParams {
		ud := l.CheckUserData(index)
		if v, ok := ud.Value.(*appParams); ok {
			return v
		}

		l.ArgError(1, "app expected")
		return nil
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

		apps := []*appParams{}
		for i := 1; i <= amount; i++ {
			app := checkAppFn(l, i)
			if app == nil {
				return 1
			}

			apps = append(apps, app)
		}

		// add the chain params to the XChain:
		app.chain = &chainParams{
			apps: apps,
		}

		// set the value:
		ud.Value = app.chain

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
