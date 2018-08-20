package modules

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestKeys_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the module:
	createXMN(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/xkeys_test.lua")

}

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
