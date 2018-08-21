package modules

import (
	"encoding/gob"

	objects "github.com/XMNBlockchain/datamint/objects"
	lua "github.com/yuin/gopher-lua"
)

const luaTables = "xtables"

// XTables represents the xtables instance
type XTables struct {
	l    *lua.LState
	objs objects.Objects
}

// CreateXTables creates a new XTables instance:
func CreateXTables(l *lua.LState) *XTables {
	// create the instance:
	out := XTables{
		l:    l,
		objs: objects.SDKFunc.Create(),
	}

	//registers the xtables module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XTables) register() {
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
		ud.Value = app.objs
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

	mt := app.l.NewTypeMetatable(luaTables)
	app.l.SetGlobal(luaTables, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadFn))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}
