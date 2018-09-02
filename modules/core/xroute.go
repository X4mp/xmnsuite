package core

import (
	applications "github.com/XMNBlockchain/datamint/applications"
	lua "github.com/yuin/gopher-lua"
)

const luaRoute = "xroute"

type routeParams struct {
	pattern  string
	saveTrx  *lua.LFunction
	delTrx   *lua.LFunction
	queryTrx *lua.LFunction
}

// XRoute represents the xroute instance
type XRoute struct {
	l           *lua.LState
	typeMapping map[string]int
	rtes        []*routeParams
}

// CreateXRoute creates a new XRoute instance:
func CreateXRoute(l *lua.LState) *XRoute {
	// create the instance:
	out := XRoute{
		l:    l,
		rtes: []*routeParams{},
		typeMapping: map[string]int{
			"retrieve": applications.Retrieve,
			"save":     applications.Save,
			"delete":   applications.Delete,
		},
	}

	//registers the xcrypto module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XRoute) register() {
	// create a new route instance:
	newRoute := func(l *lua.LState) int {
		ud := l.NewUserData()

		amount := l.GetTop()
		if amount != 3 {
			l.ArgError(1, "the exists func expected 3 parameters")
			return 1
		}

		rteTypeAsString := l.CheckString(1)
		if _, ok := app.typeMapping[rteTypeAsString]; !ok {
			l.ArgError(2, "the passed route type is invalid")
			return 1
		}

		patternAsString := l.CheckString(2)
		luaHandlr := l.CheckFunction(3)

		newRte := routeParams{
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
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaRoute)
	app.l.SetGlobal(luaRoute, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newRoute))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
