package core

import (
	lua "github.com/yuin/gopher-lua"
)

const luaRouter = "xrouter"

type routerParams struct {
	key  string
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
	// convert the table argument to a router:
	fromTableToRouterFn := func(l *lua.LState) *routerParams {
		tb := l.ToTable(1)

		routes := []*routeParams{}
		key := tb.RawGet(lua.LString("key"))
		if rawRoutes, ok := tb.RawGet(lua.LString("routes")).(*lua.LTable); ok {
			rawRoutes.ForEach(func(key lua.LValue, rawRoute lua.LValue) {
				if oneRouteUD, ok := rawRoute.(*lua.LUserData); ok {
					if oneRoute, ok := oneRouteUD.Value.(*routeParams); ok {
						routes = append(routes, oneRoute)
					}

				}
			})

		}

		return &routerParams{
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
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaRouter)
	app.l.SetGlobal(luaRouter, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newRouter))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
