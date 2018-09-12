package crypto

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestPrivKey_Success(t *testing.T) {

	// variables:
	scriptPath := "lua/privkey_test.lua"

	//create lua state:
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer context.Close()

	// create module:
	createModule(context)

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		t.Errorf("the returned error was expected to be nil, error retrned: %s", doFileErr.Error())
		return
	}
}
