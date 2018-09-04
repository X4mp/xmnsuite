package core

import (
	lua "github.com/yuin/gopher-lua"
)

const luaApplication = "xapp"

type appParams struct {
	version      string
	beginIndex   int
	endIndex     int
	routerParams *routerParams
}

// XApp represents the xapp instance
type XApp struct {
	l *lua.LState
}

// CreateXApp creates a new XApp instance:
func CreateXApp(l *lua.LState) *XApp {
	// create the instance:
	out := XApp{
		l: l,
	}

	//registers the xcrypto module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XApp) register() {
	//verifies that the given type is a Route instance:
	checkRouterFn := func(l *lua.LState, index int) *routerParams {
		ud := l.CheckUserData(index)
		if v, ok := ud.Value.(*routerParams); ok {
			return v
		}

		l.ArgError(1, "router expected")
		return nil
	}

	// create a new app instance:
	newApp := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount != 4 {
			l.ArgError(1, "the new function was expected to have 4 parameters")
			return 1
		}

		version := l.CheckString(1)
		beginIndex := l.CheckInt(2)
		endIndex := l.CheckInt(3)
		rter := checkRouterFn(l, 4)
		if rter == nil {
			return 1
		}

		// set the value:
		ud.Value = &appParams{
			version:      version,
			beginIndex:   beginIndex,
			endIndex:     endIndex,
			routerParams: rter,
		}

		l.SetMetatable(ud, l.GetTypeMetatable(luaApplication))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaApplication)
	app.l.SetGlobal(luaApplication, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newApp))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
