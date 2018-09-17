package chain

import (
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	datastore_module "github.com/xmnservices/xmnsuite/modules/datastore"
	json_module "github.com/xmnservices/xmnsuite/modules/json"
	lua "github.com/yuin/gopher-lua"
)

func TestModule_Success(t *testing.T) {
	// variables:
	dbPath := "./test_files"
	instanceID := uuid.NewV4()
	nodePK := ed25519.GenPrivKey()
	rootPubKeys := []crypto.PublicKey{}
	scriptPath := "tests/lua/chain.lua"
	ds := datastore.SDKFunc.Create()
	defer func() {
		os.RemoveAll(dbPath)
	}()

	//create lua state:
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer context.Close()

	// create json:
	json_module.SDKFunc.Create(json_module.CreateParams{
		Context: context,
	})

	// create the datastore module:
	dsMod := datastore_module.SDKFunc.Create(datastore_module.CreateParams{
		Context:   context,
		Datastore: ds,
	})

	// create module:
	module := createModule(context, dbPath, &instanceID, rootPubKeys, nodePK, dsMod)

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", doFileErr.Error())
		return
	}

	// spawn:
	_, nodeErr := module.Spawn()
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error retrned: %s", nodeErr.Error())
		return
	}
}
