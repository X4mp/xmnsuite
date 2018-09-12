package datastore

import (
	"encoding/gob"

	crypto "github.com/tendermint/tendermint/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/roles"
	"github.com/xmnservices/xmnsuite/users"
	lua "github.com/yuin/gopher-lua"
)

const luaTables = "tables"
const luaUsers = "users"
const luaRoles = "roles"
const luaKey = "keys"

type module struct {
	context *lua.LState
	ds      datastore.DataStore
	tables  objects.Objects
	usrs    users.Users
	rols    roles.Roles
	k       keys.Keys
}

func createModule(context *lua.LState, ds datastore.DataStore) Datastore {
	out := module{
		context: context,
		ds:      ds,
		tables:  ds.Objects(),
		usrs:    ds.Users(),
		rols:    ds.Roles(),
		k:       ds.Keys(),
	}

	out.register()

	return &out
}

func (app *module) register() {
	// preload datastore:
	app.context.PreloadModule("datastore", func(context *lua.LState) int {
		app.registerTables(context)
		app.registerUsers(context)
		app.registerRoles(context)
		app.registerKeys(context)
		return 1
	})
}

func (app *module) registerTables(context *lua.LState) int {
	//gob register:
	gob.Register(map[string]interface{}{})

	//verifies that the given type is an object instance:
	checkFn := func(l *lua.LState) objects.Objects {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(objects.Objects); ok {
			return v
		}

		l.ArgError(1, "tables expected")
		return nil
	}

	// load the Objects instance:
	loadFn := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.tables
		l.SetMetatable(ud, l.GetTypeMetatable(luaTables))
		l.Push(ud)
		return 1
	}

	//execute the save command on the objects instance:
	saveFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the save func expected at least 2 parameters")
			return 1
		}

		params := []*objects.ObjInKey{}
		for i := 2; i <= amount; i++ {
			oneObjInKey := objects.ObjInKey{}
			oneParam := l.CheckTable(i)
			oneParam.ForEach(func(name lua.LValue, value lua.LValue) {
				valueType := value.Type()
				nameAsString := name.String()
				if nameAsString == "table" && valueType == lua.LTTable {
					oneObjInKey.Obj = convertLTableToHashMap(value.(*lua.LTable))
				}

				if nameAsString == "key" && valueType == lua.LTString {
					oneObjInKey.Key = value.String()
				}
			})

			params = append(params, &oneObjInKey)
		}

		amountSaved := p.Save(params...)
		l.Push(lua.LNumber(amountSaved))
		return 1
	}

	//execute the retrieve command on the objects instance:
	retrieveFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		objInKey := objects.ObjInKey{
			Key: l.CheckString(2),
			Obj: new(map[string]interface{}),
		}

		p.Retrieve(&objInKey)
		mapResult := objInKey.Obj.(*map[string]interface{})
		ltable := convertHashMapToLTable(*mapResult)
		if ltable == nil {
			l.Push(lua.LNil)
			return 1
		}

		l.Push(ltable)
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return lenFn(p.Keys())(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkFn(l)
			return existsFn(p.Keys())(l)
		},
		"retrieve": retrieveFn,
		"search": func(l *lua.LState) int {
			p := checkFn(l)
			return searchFn(p.Keys())(l)
		},
		"save": saveFn,
		"delete": func(l *lua.LState) int {
			p := checkFn(l)
			return delFn(p.Keys())(l)
		},
	}

	mt := context.NewTypeMetatable(luaTables)
	context.SetGlobal(luaTables, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))

	// return:
	return 1
}

func (app *module) registerUsers(context *lua.LState) {
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

	mt := context.NewTypeMetatable(luaUsers)
	context.SetGlobal(luaUsers, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *module) registerRoles(context *lua.LState) {
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

	mt := context.NewTypeMetatable(luaRoles)
	context.SetGlobal(luaRoles, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadFn))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

func (app *module) registerKeys(context *lua.LState) {
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

	mt := context.NewTypeMetatable(luaKey)
	context.SetGlobal(luaKey, mt)

	// static attributes
	context.SetField(mt, "load", context.NewFunction(loadKeys))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}

// Get returns the datastore
func (app *module) Get() datastore.DataStore {
	return app.ds
}

// Replace replaces the datastore
func (app *module) Replace(newDS datastore.DataStore) {
	app.tables = newDS.Objects()
}
