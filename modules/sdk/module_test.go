package sdk

import (
	"encoding/hex"
	"os"
	"testing"

	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
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

	privKey := new(ed25519.PrivKeyEd25519)
	privKeyAsBytes, _ := hex.DecodeString("a328891040d6a18f2777b3638a3f56707a73505a3f4ba498c9fc80b962d475d821f4d96ca016b8c5a33fbc0e1fbaa11d16b9b008ade72ba4cef520d3b1edab70d70cf8f4ac")
	cdc.UnmarshalBinaryBare(privKeyAsBytes, privKey)
	rootPubKeys := []crypto.PubKey{
		privKey.PubKey(),
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
