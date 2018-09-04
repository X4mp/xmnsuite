package modules

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func callLuaFunc(context *lua.LState) {
	//execute the script:
	doFileErr := context.DoFile("script.lua")
	if doFileErr != nil {
		panic(doFileErr)
	}

	luaParams := lua.P{
		Fn:      context.GetGlobal("callThis"),
		NRet:    1,
		Protect: true,
	}

	// call the func:
	callErr := context.CallByParam(luaParams, lua.LNumber(12), lua.LNumber(21))
	if callErr != nil {
		panic(callErr)
	}

	// retrieve the returned value:
	ret := context.Get(-1)

	// remove the receives value:
	context.Pop(1)

	fmt.Printf("\n\n->->%d\n", ret)
}
