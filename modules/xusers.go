package modules

import (
	users "github.com/XMNBlockchain/datamint/users"
	lua "github.com/yuin/gopher-lua"
)

const luaUsers = "xusers"

// XUsers represents the xusers instance
type XUsers struct {
	l    *lua.LState
	usrs users.Users
}

// CreateXUsers creates a new XUsers instance:
func CreateXUsers(l *lua.LState) *XUsers {
	// create the instance:
	out := XUsers{
		l:    l,
		usrs: users.SDKFunc.Create(),
	}

	//registers the xusers module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XUsers) register() {
	//verifies that the given type is a Users instance:
	checkFn := func(l *lua.LState) users.Users {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(users.Users); ok {
			return v
		}

		l.ArgError(1, "users expected")
		return nil
	}

	// load the Users instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.usrs
		l.SetMetatable(ud, l.GetTypeMetatable(luaUsers))
		l.Push(ud)
		return 1
	}

	//execute the key command on the objects instance:
	keyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		key := p.Key(pubKey)
		l.Push(lua.LString(key))
		return 1
	}

	//execute the exists command on the objects instance:
	existsFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		exists := p.Exists(pubKey)
		l.Push(lua.LBool(exists))
		return 1
	}

	//execute the insert command on the objects instance:
	insertFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isInserted := p.Insert(pubKey)
		l.Push(lua.LBool(isInserted))
		return 1
	}

	//execute the delete command on the objects instance:
	deleteFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isDeleted := p.Delete(pubKey)
		l.Push(lua.LBool(isDeleted))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p.Objects().Keys())(l)
		},
		"key":    keyFn,
		"exists": existsFn,
		"insert": insertFn,
		"delete": deleteFn,
	}

	mt := app.l.NewTypeMetatable(luaUsers)
	app.l.SetGlobal(luaUsers, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadFn))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
