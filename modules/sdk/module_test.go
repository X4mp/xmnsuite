package sdk

import (
	"os"
	"testing"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	tests_chain_module "github.com/xmnservices/xmnsuite/modules/chain/tests"
	crypto_module "github.com/xmnservices/xmnsuite/modules/crypto"
	json_module "github.com/xmnservices/xmnsuite/modules/json"
	lua "github.com/yuin/gopher-lua"
)

func TestSDK_Success(t *testing.T) {

	// variables:
	blkchainDdbPath := "./test_files"
	blkchainScriptPath := "lua/chain.lua"
	defer func() {
		os.RemoveAll(blkchainDdbPath)
	}()

	privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: "60699820e4ba230a2590255678f1ce701f3943e2f7388854ae0af169b9871c0a",
	})

	rootPubKeys := []crypto.PublicKey{
		privKey.PublicKey(),
	}

	//create lua state:
	blkchainContext := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer blkchainContext.Close()

	// execute the blockchain script:
	node, nodeErr := tests_chain_module.ExecuteForTests(blkchainContext, blkchainScriptPath, blkchainDdbPath, rootPubKeys)
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}
	defer node.Stop()
	node.Start()

	// get the client:
	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// variables:
	scriptPath := "lua/sdk_test.lua"

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
	crypto_module.SDKFunc.Create(crypto_module.CreateParams{
		Context: context,
	})

	// create module:
	createModule(context, client)

	// execute:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", doFileErr.Error())
		return
	}
}
