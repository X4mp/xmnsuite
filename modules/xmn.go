package modules

import (
	keys "github.com/XMNBlockchain/datamint/keys"
	lua "github.com/yuin/gopher-lua"
)

const luaKey = "xkeys"
const luaLists = "xlists"
const luaSets = "xsets"
const luaObjs = "xobjects"
const luaUsers = "xusers"
const luaRoles = "xroles"

// XMN represents the XMN module:
type XMN struct {
	l *lua.LState
	k keys.Keys
}

func createXMN(l *lua.LState) *XMN {

	//create the instance:
	out := XMN{
		l: l,
		k: keys.SDKFunc.Create(),
	}

	//register the module on the lua state:
	out.registerXMN()

	//return the instance:
	return &out
}

func (app *XMN) registerXMN() {
	app.registerXMNKeys()
}

func (app *XMN) registerXMNKeys() {
	//verifies that the given type is a keys instance:
	checkKeys := func(l *lua.LState) keys.Keys {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(keys.Keys); ok {
			return v
		}

		l.ArgError(1, "keys expected")
		return nil
	}

	// create the Keys instance:
	loadKeys := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.k
		l.SetMetatable(ud, l.GetTypeMetatable(luaKey))
		l.Push(ud)
		return 1
	}

	//execute the save command on the keys instance:
	keysSave := func(l *lua.LState) int {
		p := checkKeys(l)
		if l.GetTop() == 3 {
			key := l.CheckString(2)
			value := l.CheckString(3)
			p.Save(key, value)
			return 0
		}

		l.ArgError(1, "the save func expected 2 parameters")
		return 1
	}

	// execute the retrieve command on the keys instance:
	keysRetrieve := func(l *lua.LState) int {
		p := checkKeys(l)
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

	// execute the delete command on the keys instance:
	keysDelete := func(l *lua.LState) int {
		p := checkKeys(l)
		amount := l.GetTop()
		if amount < 1 {
			l.ArgError(1, "the retrieve func expected at least 1 parameter")
			return 1
		}

		keys := []string{}
		for i := 2; i <= amount; i++ {
			oneKey := l.CheckString(i)
			keys = append(keys, oneKey)
		}

		amountDeleted := p.Delete(keys...)
		l.Push(lua.LNumber(amountDeleted))
		return 1
	}

	// the keys methods:
	var keysMethods = map[string]lua.LGFunction{
		"save":     keysSave,
		"retrieve": keysRetrieve,
		"delete":   keysDelete,
	}

	mt := app.l.NewTypeMetatable(luaKey)
	app.l.SetGlobal(luaKey, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadKeys))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), keysMethods))
}
