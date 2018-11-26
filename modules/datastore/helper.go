package datastore

import (
	"github.com/xmnservices/xmnsuite/datastore/keys"
	lua "github.com/yuin/gopher-lua"
)

func convertHashMapToLTable(hashmap map[string]interface{}) *lua.LTable {
	if len(hashmap) <= 0 {
		return nil
	}

	out := lua.LTable{}
	for keyname, value := range hashmap {
		if subHashMap, ok := value.(map[string]interface{}); ok {
			subLTable := convertHashMapToLTable(subHashMap)
			out.RawSet(lua.LString(keyname), subLTable)
			continue
		}

		out.RawSet(lua.LString(keyname), lua.LString(value.(string)))
	}

	return &out
}

func convertLTableToHashMap(table *lua.LTable) map[string]interface{} {
	hashmap := map[string]interface{}{}
	table.ForEach(func(name lua.LValue, value lua.LValue) {
		if value.Type() == lua.LTTable {
			hashmap[name.String()] = convertLTableToHashMap(value.(*lua.LTable))
			return
		}

		hashmap[name.String()] = value.String()
	})

	return hashmap
}

func existsFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
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

	return fn
}

func lenFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
		if l.GetTop() == 1 {
			amount := p.Len()
			l.Push(lua.LNumber(amount))
			return 1
		}

		l.ArgError(1, "the save func expected 0 parameter")
		return 1
	}

	return fn
}

func searchFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
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

	return fn
}

func delFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
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

	return fn
}
