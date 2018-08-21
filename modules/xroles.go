package modules

import (
	roles "github.com/XMNBlockchain/datamint/roles"
	crypto "github.com/tendermint/tendermint/crypto"
	lua "github.com/yuin/gopher-lua"
)

const luaRoles = "xroles"

// XRoles represents the xroles instance
type XRoles struct {
	l    *lua.LState
	rols roles.Roles
}

// CreateXRoles creates a new XRoles instance
func CreateXRoles(l *lua.LState) *XRoles {
	//create the instance:
	out := XRoles{
		l:    l,
		rols: roles.SDKFunc.Create(),
	}

	//registers the xroles module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XRoles) register() {
	//verifies that the given type is a Roles instance:
	checkFn := func(l *lua.LState) roles.Roles {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(roles.Roles); ok {
			return v
		}

		l.ArgError(1, "roles expected, received %d")
		return nil
	}

	// load the Roles instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.rols
		l.SetMetatable(ud, l.GetTypeMetatable(luaRoles))
		l.Push(ud)
		return 1
	}

	//execute the add command on the roles instance:
	addFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Add(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the del command on the roles instance:
	delFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Del(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the enableWriteAccess command on the roles instance:
	enableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.EnableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	//execute the disableWriteAccess command on the roles instance:
	disableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.DisableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	//execute the hasWriteAccess command on the roles instance:
	hasWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		keys := []string{}
		for i := 3; i <= amount; i++ {
			keys = append(keys, l.CheckString(i))
		}

		returnedKeys := p.HasWriteAccess(key, keys...)
		table := lua.LTable{}
		for _, oneKey := range returnedKeys {
			table.Append(lua.LString(oneKey))
		}

		l.Push(&table)
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"add":                addFn,
		"del":                delFn,
		"enableWriteAccess":  enableWriteAccessFn,
		"disableWriteAccess": disableWriteAccessFn,
		"hasWriteAccess":     hasWriteAccessFn,
	}

	mt := app.l.NewTypeMetatable(luaRoles)
	app.l.SetGlobal(luaRoles, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadFn))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
