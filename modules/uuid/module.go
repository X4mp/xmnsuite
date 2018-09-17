package uuid

import (
	uuid "github.com/satori/go.uuid"
	lua "github.com/yuin/gopher-lua"
)

const luaUUID = "uuid"

type module struct {
	context *lua.LState
}

func createModule(context *lua.LState) UUID {
	out := module{
		context: context,
	}

	out.register()

	return &out
}

func (app *module) register() {
	// preload uuid:
	app.context.PreloadModule("uuid", func(context *lua.LState) int {
		methods := map[string]lua.LGFunction{
			"new": func(context *lua.LState) int {
				return app.registerNew(context)
			},
		}

		ntable := context.NewTable()
		context.SetFuncs(ntable, methods)
		context.Push(ntable)

		app.registerUUID(context)

		return 1
	})
}

func (app *module) registerNew(context *lua.LState) int {

	createOrGen := func(context *lua.LState) *uuid.UUID {
		if context.GetTop() == 1 {
			idAsString := context.CheckString(1)
			id, idErr := uuid.FromString(idAsString)
			if idErr != nil {
				context.ArgError(1, "the given uuid (%s) is not a valid uuid v4 string")
				return nil
			}

			return &id
		}

		id := uuid.NewV4()
		return &id
	}

	id := createOrGen(context)
	ud := context.NewUserData()
	ud.Value = id

	context.SetMetatable(ud, context.GetTypeMetatable(luaUUID))
	context.Push(ud)
	return 1
}

func (app *module) registerUUID(context *lua.LState) int {
	checkFn := func(l *lua.LState) *uuid.UUID {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*uuid.UUID); ok {
			return v
		}

		l.ArgError(1, "UUID expected")
		return nil
	}

	// stringFn returns the string representation of the uuid
	stringFn := func(l *lua.LState) int {
		id := checkFn(l)
		l.Push(lua.LString(id.String()))
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"string": stringFn,
	}

	mt := context.NewTypeMetatable(luaUUID)

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))

	// return:
	return 1
}
