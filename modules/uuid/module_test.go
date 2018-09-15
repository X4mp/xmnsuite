package uuid

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestUUID_Success(t *testing.T) {
	execute(t, "lua/uuid_test.lua")
}

func execute(t *testing.T, scriptPath string) {
	//create lua state:
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer context.Close()

	// create module:
	createModule(context)

	//execute:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", doFileErr.Error())
		return
	}
}
