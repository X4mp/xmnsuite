package core

import (
	keys "github.com/XMNBlockchain/datamint/keys"
	lua "github.com/yuin/gopher-lua"
)

const luaKey = "xkeys"

// XKeys represents the xkeys instance
type XKeys struct {
	l *lua.LState
	k keys.Keys
}

// CreateXKeys creates a new XKeys instance:
func CreateXKeys(l *lua.LState) *XKeys {
	// create the instance:
	out := XKeys{
		l: l,
		k: keys.SDKFunc.Create(),
	}

	//registers the xkeys module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XKeys) register() {
	//verifies that the given type is a keys instance:
	checkFn := func(l *lua.LState) keys.Keys {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(keys.Keys); ok {
			return v
		}

		l.ArgError(1, "keys expected")
		return nil
	}

	// load the Keys instance:
	loadKeys := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.k
		l.SetMetatable(ud, l.GetTypeMetatable(luaKey))
		l.Push(ud)
		return 1
	}

	// execute the retrieve command on the keys instance:
	retrieveFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount != 2 {
			l.ArgError(1, "the retrieve func expected 1 parameter")
			return 1
		}

		key := l.CheckString(2)
		value := p.Retrieve(key)
		if value == nil {
			l.Push(&lua.LNilType{})
			return 1
		}

		l.Push(lua.LString(value.(string)))
		return 1
	}

	//execute the save command on the keys instance:
	saveFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 3 {
			key := l.CheckString(2)
			value := l.CheckString(3)
			p.Save(key, value)
			return 0
		}

		l.ArgError(1, "the save func expected 2 parameters")
		return 1
	}

	// the keys methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p)(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkFn(l)
			return existsFn(p)(l)
		},
		"retrieve": retrieveFn,
		"search": func(l *lua.LState) int {
			p := checkFn(l)
			return searchFn(p)(l)
		},
		"save": saveFn,
		"delete": func(l *lua.LState) int {
			p := checkFn(l)
			return delFn(p)(l)
		},
	}

	mt := app.l.NewTypeMetatable(luaKey)
	app.l.SetGlobal(luaKey, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadKeys))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
