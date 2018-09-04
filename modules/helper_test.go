package modules

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestExecute_Success(t *testing.T) {
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})

	callLuaFunc(context)

}
