package modules

import (
	"encoding/gob"

	keys "github.com/XMNBlockchain/datamint/keys"
	objects "github.com/XMNBlockchain/datamint/objects"
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
	l   *lua.LState
	k   keys.Keys
	obj objects.Objects
}

func createXMN(l *lua.LState) *XMN {

	//create the instance:
	out := XMN{
		l:   l,
		k:   keys.SDKFunc.Create(),
		obj: objects.SDKFunc.Create(),
	}

	//register the module on the lua state:
	out.register()

	//return the instance:
	return &out
}

func (app *XMN) register() {
	app.registerKeys()
	app.registerObjects()
}

func (app *XMN) registerKeys() {
	//verifies that the given type is a keys instance:
	checkKeys := func(l *lua.LState) keys.Keys {
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

	//execute the len command on the keys instance:
	keysLen := func(l *lua.LState) int {
		p := checkKeys(l)
		if l.GetTop() == 1 {
			amount := p.Len()
			l.Push(lua.LNumber(amount))
			return 1
		}

		l.ArgError(1, "the save func expected 0 parameter")
		return 1
	}

	//execute the exists command on the keys instance:
	keysExists := func(l *lua.LState) int {
		p := checkKeys(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the save func expected 0 parameter")
			return 1
		}

		keys := []string{}
		for i := 2; i <= amount; i++ {
			oneKey := l.CheckString(i)
			keys = append(keys, oneKey)
		}

		existsAmount := p.Exists(keys...)
		l.Push(lua.LNumber(existsAmount))
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

	// execute the retrieve command on the keys instance:
	keysSearch := func(l *lua.LState) int {
		p := checkKeys(l)
		amount := l.GetTop()
		if amount != 2 {
			l.ArgError(1, "the retrieve func expected 1 parameter")
			return 1
		}

		pattern := l.CheckString(2)
		results := p.Search(pattern)

		keys := lua.LTable{}
		for index, oneResult := range results {
			keys.Insert(index, lua.LString(oneResult))
		}

		l.Push(&keys)
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
		"len":      keysLen,
		"exists":   keysExists,
		"retrieve": keysRetrieve,
		"search":   keysSearch,
		"save":     keysSave,
		"delete":   keysDelete,
	}

	mt := app.l.NewTypeMetatable(luaKey)
	app.l.SetGlobal(luaKey, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadKeys))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), keysMethods))
}

func (app *XMN) registerObjects() {

	//gob register:
	gob.Register(map[string]interface{}{})

	//verifies that the given type is a keys instance:
	checkObjects := func(l *lua.LState) objects.Objects {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(objects.Objects); ok {
			return v
		}

		l.ArgError(1, "objects expected")
		return nil
	}

	// load the Objects instance:
	loadObjects := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.obj
		l.SetMetatable(ud, l.GetTypeMetatable(luaObjs))
		l.Push(ud)
		return 1
	}

	//execute the save command on the objects instance:
	saveFn := func(l *lua.LState) int {
		p := checkObjects(l)
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
				if nameAsString == "object" && valueType == lua.LTTable {
					oneObjInKey.Obj = app.convertLTableToHashMap(value.(*lua.LTable))
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
		p := checkObjects(l)
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
		ltable := app.convertHashMapToLTable(*mapResult)
		l.Push(ltable)
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"save":     saveFn,
		"retrieve": retrieveFn,
	}

	mt := app.l.NewTypeMetatable(luaObjs)
	app.l.SetGlobal(luaObjs, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadObjects))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}

func (app *XMN) convertHashMapToLTable(hashmap map[string]interface{}) *lua.LTable {
	out := lua.LTable{}
	for keyname, value := range hashmap {
		if subHashMap, ok := value.(map[string]interface{}); ok {
			subLTable := app.convertHashMapToLTable(subHashMap)
			out.RawSet(lua.LString(keyname), subLTable)
			continue
		}

		out.RawSet(lua.LString(keyname), lua.LString(value.(string)))
	}

	return &out
}

func (app *XMN) convertLTableToHashMap(table *lua.LTable) map[string]interface{} {
	hashmap := map[string]interface{}{}
	table.ForEach(func(name lua.LValue, value lua.LValue) {
		if value.Type() == lua.LTTable {
			hashmap[name.String()] = app.convertLTableToHashMap(value.(*lua.LTable))
			return
		}

		hashmap[name.String()] = value.String()
	})

	return hashmap
}
