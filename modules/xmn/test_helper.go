package xmn

import (
	lua "github.com/yuin/gopher-lua"
)

func createLuaState() *lua.LState {
	l := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})

	return l
}

func executeChunkForTests(l *lua.LState, filePath string) {
	err := l.DoFile(filePath)
	if err != nil {
		panic(err)
	}
}
