package core

import (
	"errors"
	"fmt"
	"strconv"

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
	// convert the table argument to an application:
	fromTableToAppFn := func(l *lua.LState) (*appParams, error) {
		tb := l.ToTable(1)
		version := tb.RawGet(lua.LString("version"))
		beginIndex := tb.RawGet(lua.LString("beginBlockIndex"))
		endIndex := tb.RawGet(lua.LString("endBlockIndex"))
		rterTable := tb.RawGet(lua.LString("router")).(*lua.LUserData)
		if router, ok := rterTable.Value.(*routerParams); ok {

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

			return &appParams{
				version:      version.String(),
				beginIndex:   beginIndexAsInt,
				endIndex:     endIndexAsInt,
				routerParams: router,
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
	var methods = map[string]lua.LGFunction{}

	mt := app.l.NewTypeMetatable(luaApplication)
	app.l.SetGlobal(luaApplication, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newApp))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
