package core

import (
	lua "github.com/yuin/gopher-lua"
)

const luaRouter = "xrouter"

type routerParams struct {
	rtes []*routeParams
}

// XRouter represents the xrouter instance
type XRouter struct {
	l *lua.LState
}

// CreateXRouter creates a new XRouter instance:
func CreateXRouter(l *lua.LState) *XRouter {
	// create the instance:
	out := XRouter{
		l: l,
	}

	//registers the xcrypto module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XRouter) register() {
	//verifies that the given type is a Route instance:
	checkRouteFn := func(l *lua.LState, index int) *routeParams {
		ud := l.CheckUserData(index)
		if v, ok := ud.Value.(*routeParams); ok {
			return v
		}

		l.ArgError(1, "route expected")
		return nil
	}

	// create a new router instance:
	newRouter := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount < 1 {
			l.ArgError(1, "the new function was expected to have at least 1 parameter")
			return 1
		}

		rtes := []*routeParams{}
		for i := 1; i < amount; i++ {
			oneRte := checkRouteFn(l, i)
			rtes = append(rtes, oneRte)
		}

		// set the value:
		ud.Value = &routerParams{
			rtes,
		}

		l.SetMetatable(ud, l.GetTypeMetatable(luaRouter))
		l.Push(ud)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaRouter)
	app.l.SetGlobal(luaRouter, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newRouter))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
