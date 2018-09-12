package datastore

import (
	"testing"

	datastore "github.com/xmnservices/xmnsuite/datastore"
	module_crypto "github.com/xmnservices/xmnsuite/modules/crypto"
	json_module "github.com/xmnservices/xmnsuite/modules/json"
	lua "github.com/yuin/gopher-lua"
)

func TestKeys_Success(t *testing.T) {
	execute(t, "lua/keys_test.lua")
}

func TestRoles_Success(t *testing.T) {
	execute(t, "lua/roles_test.lua")
}

func TestTables_Success(t *testing.T) {
	execute(t, "lua/tables_test.lua")
}

func TestUsers_Success(t *testing.T) {
	execute(t, "lua/users_test.lua")
}

func execute(t *testing.T, scriptPath string) {
	// variables:
	ds := datastore.SDKFunc.Create()

	//create lua state:
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer context.Close()

	// preload JSON:
	json_module.SDKFunc.Create(json_module.CreateParams{
		Context: context,
	})

	// preload crypto:
	module_crypto.SDKFunc.Create(module_crypto.CreateParams{
		Context: context,
	})

	// create module:
	createModule(context, ds)

	//execute:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", doFileErr.Error())
		return
	}
}
